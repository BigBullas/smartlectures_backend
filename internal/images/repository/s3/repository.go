package s3

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/google/uuid"
	"github.com/t1d333/smartlectures/internal/images"
	"github.com/t1d333/smartlectures/internal/images/repository"
	"github.com/t1d333/smartlectures/pkg/logger"
	"golang.org/x/net/context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Repository struct {
	path   string
	bucket string
	url    string
	logger logger.Logger
	client *minio.Client
}

func (r *Repository) UploadImage(img io.Reader, ctx context.Context) (string, error) {
	id := uuid.NewString()
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(r.path + id),
		ACL:    types.ObjectCannedACLPublicRead,
		Body:   img,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to s3 bucket: %w", err)
	}

	return fmt.Sprintf("%s/%s/%s%s", r.url, r.bucket, r.path, id), nil
}

func NewRepository(logger logger.Logger, appCfg images.Config) (repository.Repository, error) {
	// customResolver := aws.EndpointResolverWithOptionsFunc(
	// 	func(service, region string, options ...interface{}) (aws.Endpoint, error) {
	// 		if service == s3.ServiceID && region == appCfg.Region {
	// 			return aws.Endpoint{
	// 				URL:           appCfg.URL,
	// 				SigningRegion: appCfg.Region,
	// 			}, nil
	// 		}
	// 		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	// 	},
	// )

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load s3 config: %w", err)
	}

	// client := s3.NewFromConfig(cfg)

	client, err := minio.New(minioCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.accessKey, cfg.secretKey, ""),
		Secure: cfg.useSSL,
	})
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to connect to Minio: %v", err))
	}

	return &Repository{
		logger: logger,
		client: client,
		url:    appCfg.URL,
		bucket: cfg.bucket,
		path:   "attachments/",
	}, nil
}
