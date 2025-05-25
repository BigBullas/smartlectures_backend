package s3

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/t1d333/smartlectures/internal/images"
	"github.com/t1d333/smartlectures/internal/images/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
	"golang.org/x/net/context"
)

type Repository struct {
	bucket string
	url    string
	logger logger.Logger
	client *s3.Client
}

func (r *Repository) UploadImage(img io.Reader, ctx context.Context) (string, error) {
	id := uuid.NewString()
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(id),
		Body:   img,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to MinIO bucket: %w", err)
	}

	return fmt.Sprintf("%s/%s/%s", r.url, r.bucket, id), nil
}

func NewRepository(logger logger.Logger, appCfg images.Config) (repository.Repository, error) {
	// Создаём кастомный resolver для MinIO
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					URL:           appCfg.URL,    // http://localhost:9000
					SigningRegion: appCfg.Region, // ru-msk (можно любое значение)
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		},
	)

	// Конфигурация с использованием статических credentials (логин/пароль MinIO)
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				appCfg.AccessKey, // minioadmin (или ваш ключ)
				appCfg.SecretKey, // minioadmin (или ваш секрет)
				"",               // Сессионный токен (не нужен для MinIO)
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load MinIO config: %w", err)
	}

	// Отключаем SSL, если используем http://
	if !appCfg.UseSSL {
		cfg.HTTPClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	client := s3.NewFromConfig(cfg)
	return &Repository{
		logger: logger,
		client: client,
		url:    appCfg.URL,
		bucket: appCfg.BucketName,
	}, nil
}
