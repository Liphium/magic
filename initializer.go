package magic

import (
	"log"

	"github.com/Liphium/magic/mconfig"
	"github.com/spf13/pflag"
)

type Config struct {

	// Required. This app name will be used in the container name of databases or other services we may start for you automatically. Make sure you have no other project using the same name.
	AppName string

	// Required. This function will be executed to plan all the containers you want to start. You may create databases and more using the context passed into this function.
	PlanDeployment func(ctx *mconfig.Context)

	// Required. This should start your app like normal. Expect all the database and other containers to be started at this point. Magic will make sure they're all ready by doing health checks and stuff.
	StartFunction func()
}

// Start Magic. Make sure to provide all the arguments in the config that are required.
func Start(config Config) {
	if config.AppName == "" {
		log.Fatalln("You MUST provide an app name for Magic's config. It will be used in the container name of databases or other services we may start for you automatically. Make sure no two app names are the same as containers might otherwise be deleted without you expecting it.")
	}

	// Enable verbose logging in case desired
	mconfig.VerboseLogging = *pflag.Bool("m-verbose", false, "Enable verbose logging for Magic.")

	// Create a new Magic context
	currentProfile := *pflag.String("m-profile", "default", "The profile to use for Magic.")
	ctx := mconfig.DefaultContext(config.AppName, currentProfile)

	// Plan the deployment
	config.PlanDeployment(ctx)
}
