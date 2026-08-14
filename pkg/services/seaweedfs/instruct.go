package seaweedfs

import (
	"context"
	"fmt"

	"github.com/Liphium/magic/v4/mconfig"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/moby/moby/client"
)

// Handles instructions for SeaweedFS.
//
// This will only handle the clear files instruction.
func (sd *SeaweedFSDriver) HandleInstruction(ctx context.Context, c *client.Client, container mconfig.ContainerInformation, instruction mconfig.Instruction) error {
	switch instruction {
	case mconfig.InstructionClearFiles:
		return sd.clearFiles(container)
	}
	return nil
}

func (sd *SeaweedFSDriver) clearFiles(container mconfig.ContainerInformation) error {
	client, err := sd.getS3Client(container.Ports[0])
	if err != nil {
		return fmt.Errorf("s3: %w", err)
	}

	// Delete all of the buckets the driver created previously
	for _, bucket := range sd.Buckets {
		_, err = client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("couldn't delete bucket %s: %w", bucket, err)
		}
	}

	return nil
}
