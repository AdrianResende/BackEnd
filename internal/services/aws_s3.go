package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	Client     *s3.Client
	BucketName string
	Region     string
}

func NewS3Service() (*S3Service, error) {
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_BUCKET_NAME")

	if region == "" || bucket == "" {
		return nil, fmt.Errorf("AWS_REGION ou AWS_BUCKET_NAME não definidos")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Service{
		Client:     client,
		BucketName: bucket,
		Region:     region,
	}, nil
}

func (s *S3Service) UploadFile(file multipart.File, fileName string, contentType string) (string, error) {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return "", err
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
	}

	_, err = s.Client.PutObject(context.TODO(), input)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.BucketName, s.Region, fileName)
	return url, nil
}
