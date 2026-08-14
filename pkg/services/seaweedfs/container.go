package seaweedfs

import (
	"context"
	"fmt"

	"github.com/Liphium/magic/v4/mconfig"
	mservices "github.com/Liphium/magic/v4/mrunner/services"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/moby/moby/client"
)

// Should create a new container for the service or use the existing one
// (returns container id + error in case one happened).
func (sd *SeaweedFSDriver) CreateContainer(ctx context.Context, c *client.Client, a mconfig.ContainerAllocation) (string, error) {
	if sd.Image == "" {
		return "", fmt.Errorf("please specify a proper image")
	}

	return mservices.CreateContainer(ctx, seaweedLog, c, a, mservices.ManagedContainerOptions{
		Image: sd.Image,
		Env: []string{
			fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", SeaweedFSS3AccessKey),
			fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", SeaweedFSS3SecretKey),
		},
		Ports: RequiredPorts,
		Cmd:   []string{"server", "-s3"},
		Volumes: []mservices.ContainerVolume{
			{NameSuffix: "data", Target: "/data"},
		},
	})
}

// Check SeaweedFS health by listing the buckets inside the service (with limit of 1).
func (sd *SeaweedFSDriver) IsHealthy(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) (bool, error) {
	if mconfig.VerboseLogging {
		seaweedLog.Println("checking health...")
	}

	client, err := sd.getS3Client(container.Ports[0])
	if err != nil {
		return false, fmt.Errorf("s3: %w", err)
	}

	// List all of the buckets, if there is an error, we know the service must not be healthy yet.
	_, err = client.ListBuckets(context.Background(), &s3.ListBucketsInput{
		MaxBuckets: aws.Int32(1),
	})
	if mconfig.VerboseLogging {
		seaweedLog.Printf("health check result: %v (error: %v)", err == nil, err)
	}
	return err == nil, nil
}

// getS3Client creates a new S3 client for the bucket, this is for internal use within the driver to create buckets, etc.
func (sd *SeaweedFSDriver) getS3Client(port uint) (*s3.Client, error) {

	// Build base config with static credentials
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithRetryMaxAttempts(0),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(SeaweedFSS3AccessKey, SeaweedFSS3SecretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// Custom endpoint resolver → points at the S3-compatible store
	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("http://127.0.0.1:%d", port))
		o.UsePathStyle = true
	}), nil
}

// Initialize the container. Creates all of the buckets desired by the user.
func (sd *SeaweedFSDriver) Initialize(ctx context.Context, c *client.Client, container mconfig.ContainerInformation) error {
	client, err := sd.getS3Client(container.Ports[0])
	if err != nil {
		return fmt.Errorf("s3: %w", err)
	}

	// Try to create all of the buckets if they are not there yet
	for _, bucket := range sd.Buckets {
		_, err = client.HeadBucket(context.Background(), &s3.HeadBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			_, makeErr := client.CreateBucket(context.Background(), &s3.CreateBucketInput{
				Bucket: aws.String(bucket),
			})
			if makeErr != nil {
				return fmt.Errorf("create bucket %s: %w", bucket, makeErr)
			}
		}
	}
	return nil
}
