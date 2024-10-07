package main

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/shion13/interview.devops/internal/server"
)

func main() {
	// Using zap SugaredLogger. It's 4-10x faster than other structured logging packages and includes both structured and printf-style APIs.
	level := zap.InfoLevel
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		level = zap.DebugLevel
	case "INFO":
		level = zap.InfoLevel
	case "ERROR":
		level = zap.ErrorLevel
	}

	atom := zap.NewAtomicLevel()
	atom.SetLevel(level)
	encoderCfg := zap.NewProductionEncoderConfig()
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	defer logger.Sync()

	minio_user := os.Getenv("MINIO_ROOT_USER")
	minio_password := os.Getenv("MINIO_ROOT_PASSWORD")
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	// Load the Shared AWS Configuration (~/.aws/config)
	if err != nil {
		logger.Sugar().Fatalf("unable to load default aws config %w", err)
	}

	if minio_user != "" {
		// TODO Evaluate use of endpoint resolver to coordinate AWS SDK to MINIO Server mapping
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               "minio-server:9000",
				SigningRegion:     "eu-west-2",
				HostnameImmutable: true,
			}, nil
		})

		awsConfig = aws.Config{
			Credentials:                 credentials.NewStaticCredentialsProvider(minio_user, minio_password, ""),
			EndpointResolverWithOptions: resolver,
			Region:                      "eu-west-2",
		}
	}

	s := server.Server{}
	s.Setup(awsConfig, logger)
	s.Serve()

}
