package app

import (
	"context"
	"os"
	"testing"

	"github.com/Liphium/magic/v4"
	"github.com/Liphium/magic/v4/mconfig"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

// This file contains the test setup for Magic and a simple test to verify functionality. For a real app, you should probably have a lot more tests.

func TestMain(m *testing.M) {
	magic.PrepareTesting(m, GetConfig())
}

const testFile = `Hello, this is a test file!

It has some lines of content!`

func TestFileService(t *testing.T) {
	runner := magic.GetTestRunner()
	assert.NoError(t, runner.RunInstruction(mconfig.InstructionClearFiles))

	// Prepare a local folder as a test environment
	assert.NoError(t, os.RemoveAll(".testing"))
	assert.NoError(t, os.Mkdir(".testing", os.ModePerm))
	assert.NoError(t, os.WriteFile(".testing/test.txt", []byte(testFile), os.ModePerm))

	// Upload and make sure there isn't any error
	assert.NoError(t, UploadFile(UploadData{
		UploadPath: ".testing/test.txt",
		Target:     "test.txt",
	}))

	// Check if the file really exists (we're in the same process, so we can use the same S3 client)
	_, err := Client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(Bucket),
		Key:    aws.String("test.txt"),
	})
	assert.NoError(t, err)

	// Try to download the file as well
	assert.NoError(t, DownloadFile(DownloadData{
		Path:   "test.txt",
		SaveTo: ".testing/downloaded.txt",
	}))

	// Make sure the content of the files matches
	content, err := os.ReadFile(".testing/downloaded.txt")
	assert.NoError(t, err)
	assert.Equal(t, testFile, string(content))
}
