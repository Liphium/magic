package app

import (
	"fmt"
	"time"

	"github.com/Liphium/magic/pkg/services/seaweedfs"
	"github.com/Liphium/magic/v4"
	"github.com/Liphium/magic/v4/mconfig"
	"github.com/Liphium/magic/v4/scripting"
)

// The config for Magic is defined in this file.

func GetConfig() magic.Config {
	timeout := 2 * time.Minute

	return magic.Config{
		AppName: "magic-example-file-sharing",
		PlanDeployment: func(ctx *mconfig.Context) {

			// We need the SeaweedFS driver to create a S3-compatible storage bucket
			const BucketName = "files"
			driver := seaweedfs.NewDriver("chrislusf/seaweedfs:4.41").
				NewBucket(BucketName) // Create the bucket we want

			// Register the driver so Magic knows about it
			ctx.Register(driver)

			// Allocate a port for the API
			port := ctx.ValuePort(8080) // We want it to be 8080, but Magic may allocate a different one if it's already in use

			// Add all of the environment variables required for the connection to S3
			ctx.WithEnvironment(mconfig.Environment{
				"STORAGE_HOST": mconfig.ValueWithBase([]mconfig.EnvironmentValue{driver.Host(ctx), driver.Port(ctx)}, func(s []string) string {
					return fmt.Sprintf("%s:%s", s[0], s[1])
				}),
				"STORAGE_ACCESS_KEY": driver.AccessKey(),
				"STORAGE_SECRET_KEY": driver.SecretKey(),
				"STORAGE_BUCKET":     mconfig.ValueStatic(BucketName),

				// This is the address for the web service based on the port we allocated before
				"LISTEN": mconfig.ValueWithBase([]mconfig.EnvironmentValue{port}, func(s []string) string {
					return fmt.Sprintf("127.0.0.1:%s", s[0])
				}),
			})
		},
		StartFunction: Start,

		// Scripts for uploading, downloading, and clearing files
		Scripts: []scripting.Script{
			scripting.CreateScript("upload", "Upload a file to the service", UploadFile),
			scripting.CreateScript("download", "Download a file from the service", DownloadFile),
			scripting.CreateScript("clear", "Clear all files.", ClearFiles),
		},

		// We need a little longer cause SeaweedFS takes pretty long to start up (currently 2 minutes, look at the variable above)
		TestAppTimeout: &timeout,
	}
}
