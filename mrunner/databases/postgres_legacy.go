package databases

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Liphium/magic/v2/mconfig"
	"github.com/Liphium/magic/v2/util"
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

var pgLog *log.Logger = log.New(os.Stdout, "pg-manager ", log.Default().Flags())

type PostgresDriver struct {
	image     string
	databases []string
}

// Create a new PostgreSQL legacy service driver.
//
// It currently supports version PostgreSQL v14-17. Use NewPostgresDriver for v18 and beyond.
//
// This driver will eventually be deprecated and replaced by the one for v18 and above.
func NewLegacyPostgresDriver(image string) *PostgresDriver {
	imageVersion := strings.Split(image, ":")[1]

	// Supported (confirmed and tested) major versions for this Postgres driver
	var supportedPostgresVersions = []string{"14", "15", "16", "17"}

	// Do a quick check to make sure the image version is actually supported
	supported := false
	for _, version := range supportedPostgresVersions {
		if strings.HasPrefix(imageVersion, fmt.Sprintf("%s.", version)) {
			supported = true
		}
	}
	if !supported {
		pgLog.Fatalln("ERROR: Version", imageVersion, "is currently not supported.")
	}

	return &PostgresDriver{
		image: image,
	}
}

func (pd *PostgresDriver) NewDatabase(name string) *PostgresDriver {
	pd.databases = append(pd.databases, name)
	return pd
}

// A unique identifier for the database container
func (pd *PostgresDriver) GetUniqueId() string {
	return "postgres1417"
}

func (pd *PostgresDriver) GetRequiredPortAmount() int {
	return 1
}

func (pd *PostgresDriver) GetImage() string {
	return pd.image
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
		for _, container := range ctx.Plan().Containers {
			if container.Name == mconfig.PlannedContainerName(ctx.Plan(), pd) {
				return fmt.Sprintf("%d", ctx.Plan().AllocatedPorts[container.Ports[0]])
			}
		}

		util.Log.Fatalln("ERROR: Couldn't find port for PostgreSQL container in plan!")
		return "not found"
	})
}
