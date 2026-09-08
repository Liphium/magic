package surrealdb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Liphium/magic/v4/mconfig"
	mservices "github.com/Liphium/magic/v4/mrunner/services"
	"github.com/Liphium/magic/v4/util"
)

// Make sure the driver complies
var _ mconfig.ServiceDriver = &SurrealDriver{}

// IMPORTANT: Having non-static passwords would make Magic not works as the Container allocation currently does not contain service driver data.
//
// This means that instruction calling would break if we added back password and username changing.
const (
	SurrealUsername = "root"
	SurrealPassword = "root"
)

var RequiredPorts = []string{
	"8000/tcp",
}

var surrealLog *log.Logger = log.New(os.Stdout, "surrealdb ", log.Default().Flags())

// A single SurrealDB database inside of a namespace. SurrealDB always scopes databases by namespaces, so both are required to connect to one.
type SurrealDatabase struct {
	Namespace string `json:"namespace"`
	Database  string `json:"database"`
}

type SurrealDriver struct {
	Image     string            `json:"image"`
	Databases []SurrealDatabase `json:"databases"`
}

// Make sure to register the driver
func init() {
	mconfig.RegisterDriver(&SurrealDriver{})
}

// Create a new SurrealDB service driver.
//
// It currently only supports SurrealDB v3.
func NewDriver(image string) *SurrealDriver {
	imageVersion := mservices.GetImageMajorVersion(image)

	// Supported (confirmed and tested) major versions for this SurrealDB driver
	var supportedSurrealVersions = []int{3}

	supported := false
	imageMajor := mservices.GetImageMajorVersion(image)
	for _, version := range supportedSurrealVersions {
		if imageMajor == version {
			supported = true
		}
	}
	if !supported {
		surrealLog.Fatalln("ERROR: Version", imageVersion, "is currently not supported (only v3 supported).")
	}

	return &SurrealDriver{
		Image: image,
	}
}

func (sd *SurrealDriver) Load(data string) (mconfig.ServiceDriver, error) {
	var driver SurrealDriver
	if err := json.Unmarshal([]byte(data), &driver); err != nil {
		return nil, err
	}
	return &driver, nil
}

func (sd *SurrealDriver) Save() (string, error) {
	bytes, err := json.Marshal(sd)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// NewDatabase registers a new database for the driver. Since SurrealDB always scopes databases by namespaces, both the namespace and the database name have to be specified.
func (sd *SurrealDriver) NewDatabase(namespace string, database string) *SurrealDriver {
	sd.Databases = append(sd.Databases, SurrealDatabase{
		Namespace: namespace,
		Database:  database,
	})
	return sd
}

// A unique identifier for the database driver. This is appended to the container name to make sure we know it's the container from the driver.
func (sd *SurrealDriver) GetUniqueId() string {
	return "surrealdb"
}

func (sd *SurrealDriver) GetRequiredPorts() []string {
	return RequiredPorts
}

func (sd *SurrealDriver) GetImage() string {
	return sd.Image
}

// Get the username of the databases in this driver as a EnvironmentValue for your config.
func (sd *SurrealDriver) Username() mconfig.EnvironmentValue {
	return mconfig.ValueStatic(SurrealUsername)
}

// Get the password for the user of the databases in this driver as a EnvironmentValue for your config.
func (sd *SurrealDriver) Password() mconfig.EnvironmentValue {
	return mconfig.ValueStatic(SurrealPassword)
}

// Get hostname of the database container created by the driver as a EnvironmentValue for your config.
func (sd *SurrealDriver) Host(ctx *mconfig.Context) mconfig.EnvironmentValue {
	return mconfig.ValueStatic("127.0.0.1")
}

// Get the port of the database container created by the driver as a EnvironmentValue for your config.
func (sd *SurrealDriver) Port(ctx *mconfig.Context) mconfig.EnvironmentValue {
	return mconfig.ValueFunction(func() string {
		for id, container := range ctx.Plan().Containers {
			if id == sd.GetUniqueId() {
				return fmt.Sprintf("%d", ctx.Plan().AllocatedPorts[container.Ports[0]])
			}
		}

		util.Log.Fatalln("ERROR: Couldn't find port for SurrealDB container in plan!")
		return "not found"
	})
}
