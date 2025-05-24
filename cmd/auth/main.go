package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"encoding/json"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/t1d333/smartlectures/cmd/auth/config"
	"github.com/t1d333/smartlectures/internal/auth"
	grpcDel "github.com/t1d333/smartlectures/internal/auth/delivery/grpc"
	delivery "github.com/t1d333/smartlectures/internal/auth/delivery/http"
	"github.com/t1d333/smartlectures/internal/auth/repository"
	"github.com/t1d333/smartlectures/internal/auth/service"
	middl "github.com/t1d333/smartlectures/internal/middleware"
	"github.com/t1d333/smartlectures/pkg/logger"

	storage "github.com/t1d333/smartlectures/internal/storage"
)

func main() {
	logger := logger.NewLogger()

	cfg, err := config.NewConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		logger.Fatal(err)
	}

	configJSON, _ := json.MarshalIndent(cfg, "", "  ")
	logger.Infof("Full config:\n %s", string(configJSON))

	router := gin.Default()

	router.Use(middl.NewErrorHandler(logger))
	router.Use(cors.Default())

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.DragonflyURL,
		Password: cfg.DragonflyPassword,
		DB:       0,
	})

	if status := redisClient.Ping(context.Background()); status.Err() != nil {
		logger.Fatal(status.Err())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("creating storage client...")
	conn, err := grpc.Dial(
		cfg.StorageAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatalf("failed to create storage client %s", err)
	}
	defer conn.Close()

	logger.Info("storage client cratead successfully")

	storageClient := storage.NewStorageClient(conn)

	rep := repository.New(logger, redisClient, pool)
	svc := service.New(logger, rep, storageClient)
	del := delivery.NewDelivery(logger, router, svc)

	del.RegisterHandler()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			logger.Fatalf("failed to listen: %v", err)
		}

		var opts []grpc.ServerOption

		del := grpcDel.NewDelivery(logger, svc)
		logger.Info("starting grpc server on port 50051...")
		grpcServer := grpc.NewServer(opts...)
		auth.RegisterAuthServiceServer(grpcServer, del)

		if err = grpcServer.Serve(lis); err != nil {
			logger.Fatal("failed to start grpc server", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	logger.Info("timeout of 5 seconds.")
	logger.Info("Server exiting")
}
