package bucket

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3User struct {
	client *s3.Client
}

func SetupS3User(cfg aws.Config) s3User {

	return s3User{client: s3.NewFromConfig(cfg)}
}

func (s s3User) PushFileToBucket(bucket string, key string, file io.Reader) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{Bucket: &bucket, Key: &key, Body: file})
	if err != nil {
		return fmt.Errorf("unable to push object into s3 bucket %w", err)
	}
	return nil
}
