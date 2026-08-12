package seaweedfs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Liphium/magic/v3/mconfig"
	mservices "github.com/Liphium/magic/v3/mrunner/services"
	"github.com/Liphium/magic/v3/util"
)

// Make sure the driver complies with the interface.
var _ mconfig.ServiceDriver = &SeaweedFSDriver{}

const (
	ServiceName          = "seaweedfs"
	SeaweedFSS3AccessKey = "admin"
	SeaweedFSS3SecretKey = "secret"
)

var RequiredPorts = []string{
	"8333/tcp",
}

var seaweedLog *log.Logger = log.New(os.Stdout, "seaweedfs ", log.Default().Flags())

type SeaweedFSDriver struct {
	Buckets []string `json:"buckets"`
	Image   string   `json:"image"`
}

// Make sure to register the driver.
func init() {
	mconfig.RegisterDriver(&SeaweedFSDriver{})
}

// Create a new SeaweedFS service driver.
//
// It currently supports SeaweedFS major version 3 (which provides the S3 gateway).
func NewDriver(image string) *SeaweedFSDriver {

	// Supported (confirmed and tested) major versions for this SeaweedFS driver
	var supportedSeaweedFSVersions = []int{4}

	// Do a quick check to make sure the image version is actually supported
	supported := false
	imageMajor := mservices.GetImageMajorVersion(image)
	for _, version := range supportedSeaweedFSVersions {
		if imageMajor == version {
			supported = true
		}
	}
	if !supported {
		seaweedLog.Fatalln("ERROR: Major version", imageMajor, "is currently not supported.")
	}

	return &SeaweedFSDriver{
		Image: image,
	}
}

func (pd *SeaweedFSDriver) NewBucket(name string) *SeaweedFSDriver {
	pd.Buckets = append(pd.Buckets, name)
	return pd
}

func (sd *SeaweedFSDriver) Load(data string) (mconfig.ServiceDriver, error) {
	var driver SeaweedFSDriver
	if err := json.Unmarshal([]byte(data), &driver); err != nil {
		return nil, err
	}
	return &driver, nil
}

func (sd *SeaweedFSDriver) Save() (string, error) {
	bytes, err := json.Marshal(sd)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// A unique identifier for the SeaweedFS driver. Appended to the container name
// to make sure we know it's the container from the driver.
func (sd *SeaweedFSDriver) GetUniqueId() string {
	return ServiceName
}

func (sd *SeaweedFSDriver) GetRequiredPorts() []string {
	return RequiredPorts
}

func (sd *SeaweedFSDriver) GetImage() string {
	return sd.Image
}

// Host of the S3 endpoint as an EnvironmentValue for your config.
func (sd *SeaweedFSDriver) Host(ctx *mconfig.Context) mconfig.EnvironmentValue {
	return mconfig.ValueStatic("127.0.0.1")
}

// Port of the S3 gateway as an EnvironmentValue for your config.
func (sd *SeaweedFSDriver) Port(ctx *mconfig.Context) mconfig.EnvironmentValue {
	return mconfig.ValueFunction(func() string {
		for id, container := range ctx.Plan().Containers {
			if id == sd.GetUniqueId() {
				return fmt.Sprintf("%d", ctx.Plan().AllocatedPorts[container.Ports[0]])
			}
		}

		util.Log.Fatalln("ERROR: Couldn't find port for SeaweedFS container in plan!")
		return "not found"
	})
}

// The full S3 endpoint URL (host:port) as an EnvironmentValue for your config.
func (sd *SeaweedFSDriver) Endpoint(ctx *mconfig.Context) mconfig.EnvironmentValue {
	return mconfig.ValueWithBase([]mconfig.EnvironmentValue{sd.Host(ctx), sd.Port(ctx)}, func(s []string) string {
		return s[0] + ":" + s[1]
	})
}

// Access key accepted by the S3 gateway as an EnvironmentValue for your config.
func (sd *SeaweedFSDriver) AccessKey() mconfig.EnvironmentValue {
	return mconfig.ValueStatic(SeaweedFSS3AccessKey)
}

// Secret key accepted by the S3 gateway as an EnvironmentValue for your config.
func (sd *SeaweedFSDriver) SecretKey() mconfig.EnvironmentValue {
	return mconfig.ValueStatic(SeaweedFSS3SecretKey)
}
