package s3storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/DiegoJCordeiro/image-resizing-go-producer-api/internal/domain/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	bucketName = "diegojc-bucket"
)

type S3StorageImpl struct {
	bucketName string
	s3Client   *s3.Client
}

func NewS3StorageImpl(awsEndpoint string) *S3StorageImpl {
	options := []func(*config.LoadOptions) error{
		config.WithRegion("us-east-1"),
	}

	if awsEndpoint != "" {
		options = append(options, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		))
	}

	configuration, err := config.LoadDefaultConfig(context.Background(), options...)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(configuration, func(o *s3.Options) {
		o.UsePathStyle = true
		if awsEndpoint != "" {
			o.BaseEndpoint = aws.String(awsEndpoint)
		}
	})

	return &S3StorageImpl{
		bucketName: bucketName,
		s3Client:   s3Client,
	}
}

func (storageService *S3StorageImpl) GetObject(fullPath string) ([]byte, error) {

	ctx := context.Background()
	ctxTimeout, cancelFunc := context.WithTimeout(ctx, time.Second*60)

	defer cancelFunc()

	objectOutput, err := storageService.s3Client.GetObject(ctxTimeout, &s3.GetObjectInput{
		Bucket: &storageService.bucketName,
		Key:    &fullPath,
	})

	if err != nil {
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &noSuchKey) {
			return nil, fmt.Errorf("object not found: %s", fullPath)
		}
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}

	defer objectOutput.Body.Close()

	objectBody, err := io.ReadAll(objectOutput.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read object body: %w", err)
	}

	return objectBody, nil
}

func (storageService *S3StorageImpl) PutObject(imageModel *models.ImageModel) error {

	ctx := context.Background()
	ctxTimeout, cancelFunc := context.WithTimeout(ctx, time.Second*60)

	defer cancelFunc()

	contentTypeAws := http.DetectContentType(imageModel.Object)

	_, err := storageService.s3Client.PutObject(ctxTimeout, &s3.PutObjectInput{
		Bucket:        &storageService.bucketName,
		Key:           &imageModel.FullPath,
		Body:          bytes.NewReader(imageModel.Object),
		ContentLength: aws.Int64(int64(len(imageModel.Object))),
		ContentType:   aws.String(contentTypeAws),
	})

	if err != nil {
		return err
	}

	return nil
}
