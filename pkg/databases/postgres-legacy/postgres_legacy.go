package postgres_legacy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Liphium/magic/v3/mconfig"
	mservices "github.com/Liphium/magic/v3/mrunner/services"
	"github.com/Liphium/magic/v3/util"
	_ "github.com/lib/pq"
)

// Make sure the driver complies
var _ mconfig.ServiceDriver = &PostgresDriver{}

// IMPORTANT: Having non-static passwords would make Magic not works as the Container allocation currently does not contain service driver data.
//
// This means that instruction calling would break if we added back password and username changing.
const (
	PostgresUsername = "postgres"
	PostgresPassword = "postgres"
)

var pgLegacyLog *log.Logger = log.New(os.Stdout, "pg-legacy ", log.Default().Flags())

type PostgresDriver struct {
	Image     string   `json:"image"`
	Databases []string `json:"databases"`
}

// Create a new PostgreSQL legacy service driver.
//
// It currently supports version PostgreSQL 14-17, older versions have not been tested. Use the new postgres driver for PostgreSQL 18 and beyond.
//
// This driver is deprecated and will be removed when PostgreSQL 20 comes out.
func NewDriver(image string) *PostgresDriver {
	imageVersion := strings.Split(image, ":")[1]

	// Supported (confirmed and tested) major versions for this Postgres driver
	var supportedPostgresVersions = []int{14, 15, 16, 17}

	// Do a quick check to make sure the image version is actually supported
	supported := false
	imageMajor := mservices.GetImageMajorVersion(image)
	for _, version := range supportedPostgresVersions {
		if imageMajor == version {
			supported = true
		}
	}
	if !supported {
		pgLegacyLog.Fatalln("ERROR: Version", imageVersion, "is currently not supported.")
	}

	return &PostgresDriver{
		Image: image,
	}
}

func (pd *PostgresDriver) Load(data string) (mconfig.ServiceDriver, error) {
	var driver PostgresDriver
	if err := json.Unmarshal([]byte(data), &driver); err != nil {
		return nil, err
	}
	return &driver, nil
}

func (pd *PostgresDriver) Save() (string, error) {
	bytes, err := json.Marshal(pd)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (pd *PostgresDriver) NewDatabase(name string) *PostgresDriver {
	pd.Databases = append(pd.Databases, name)
	return pd
}

// A unique identifier for the database driver. This is appended to the container name to make sure we know it's the container from the driver.
func (pd *PostgresDriver) GetUniqueId() string {
	return "postgres1417" // Context for this: Since this driver supports PostgreSQL v14-v17 this just makes it easier to know when seeing the container in "docker ps" or sth
}

func (pd *PostgresDriver) GetRequiredPortAmount() int {
	return 1
}

func (pd *PostgresDriver) GetImage() string {
	return pd.Image
}

// Get the username of the databases in this driver as a EnvironmentValue for your config.
func (pd *PostgresDriver) Username() mconfig.EnvironmentValue {
	return mconfig.ValueStatic(PostgresUsername)
}

// Get the password for the user of the databases in this driver as a EnvironmentValue for your config.
func (pd *PostgresDriver) Password() mconfig.EnvironmentValue {
	return mconfig.ValueStatic(PostgresPassword)
}

// Get hostname of the database container created by the driver as a EnvironmentValue for your config.
func (pd *PostgresDriver) Host(ctx *mconfig.Context) mconfig.EnvironmentValue {
	return mconfig.ValueStatic("127.0.0.1")
}

// Get the port of the database container created by the driver as a EnvironmentValue for your config.
func (pd *PostgresDriver) Port(ctx *mconfig.Context) mconfig.EnvironmentValue {
	return mconfig.ValueFunction(func() string {
		for id, container := range ctx.Plan().Containers {
			if id == pd.GetUniqueId() {
				return fmt.Sprintf("%d", ctx.Plan().AllocatedPorts[container.Ports[0]])
			}
		}

		util.Log.Fatalln("ERROR: Couldn't find port for PostgreSQL container in plan!")
		return "not found"
	})
}
