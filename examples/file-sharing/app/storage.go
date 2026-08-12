package app

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// This file handles connection to S3.

// Client exposes the S3 client for use across services.
var Client *s3.Client = nil

// Presign exposes the S3 presign client for presigned requests.
var Presign *s3.PresignClient = nil

// Bucket name resolved from env.
var Bucket = "space"

func Connect() {
	if Client != nil {
		return
	}

	host := os.Getenv("STORAGE_HOST")
	accessKey := os.Getenv("STORAGE_ACCESS_KEY")
	secretKey := os.Getenv("STORAGE_SECRET_KEY")
	Bucket = os.Getenv("STORAGE_BUCKET")
	region := os.Getenv("STORAGE_REGION")
	if region == "" {
		region = "us-east-1"
	}
	if Bucket == "" {
		log.Fatal("STORAGE_BUCKET must be set")
	}
	if host == "" {
		log.Fatal("STORAGE_HOST must be set")
	}

	// Build base config with static credentials
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	// Custom endpoint resolver → points at the S3-compatible store
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://" + host)
		o.UsePathStyle = true
	})

	// Ensure the bucket exists
	_, err = client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(Bucket),
	})
	if err != nil {
		log.Fatal("storage bucket does not exist")
	}

	Client = client

	log.Println("object storage: connected")
}
