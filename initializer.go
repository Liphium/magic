package magic

import (
	"fmt"
	"log"
	"os"

	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/magic/mrunner"
	"github.com/Liphium/magic/scripting"
	"github.com/spf13/pflag"
)

var Log *log.Logger = log.New(os.Stdout, "magic", log.Default().Flags())

var prepared = false

type Config struct {

	// Required. This app name will be used in the container name of databases or other services we may start for you automatically. Make sure you have no other project using the same name.
	AppName string

	// Required. This function will be executed to plan all the containers you want to start. You may create databases and more using the context passed into this function.
	PlanDeployment func(ctx *mconfig.Context)

	// Required. This should start your app like normal. Expect all the database and other containers to be started at this point. Magic will make sure they're all ready by doing health checks and stuff.
	StartFunction func()

	// Add scripts that could be useful to you while developing this app.
	Scripts []scripting.Script
}

// Start your application with Magic. Make sure to provide all the arguments in the config that are required.
func Start(config Config) {
	Prepare(config, false)

	// Start the app after everything is prepared
	config.StartFunction()
}

// Start all the containers and more. Call this before running tests with Magic.
func Prepare(config Config, tests bool) {

	// Make sure we're not preparing again (this could happen in tests)
	if prepared {
		return
	}
	prepared = true

	if config.AppName == "" {
		log.Fatalln("You MUST provide an app name for Magic's config. It will be used in the container name of databases or other services we may start for you automatically. Make sure no two app names are the same as containers might otherwise be deleted without you expecting it.")
	}

	// Enable verbose logging in case desired
	mconfig.VerboseLogging = *pflag.Bool("m-verbose", false, "Enable verbose logging for Magic.")

	// Create a new Magic context
	currentProfile := *pflag.String("profile", "default", "The profile to use for Magic.")
	if tests {
		currentProfile = "test-" + currentProfile
	}
	ctx := mconfig.DefaultContext(config.AppName, currentProfile)

	// Check if all scripts should be listed
	if *pflag.Bool("scripts", false, "Lists all scripts registered in Magic.") {

		fmt.Println()
		fmt.Println("Listing all scripts registered in Magic.")
		fmt.Println("Use -r or --run <script> to run any of them.")
		fmt.Println()
		for _, script := range config.Scripts {
			fmt.Printf("%s - %s", script.Name, script.Description)
		}

		return
	}

	// Create a factory for initializing everything
	factory, err := createFactory()
	if err != nil {
		Log.Fatalln("Something went wrong:", err)
		return
	}
	factory.WarnIfNotIgnored()

	// Check if a script should be run
	script := *pflag.StringP("run", "r", "", "Provide this if you want to run a script.")
	if script != "" {
		// TODO: Actually run the script
		return
	}

	// Plan the deployment
	config.PlanDeployment(ctx)

	// Create the runner and deploy containers
	runner, err := mrunner.NewRunner(ctx)
	if err != nil {
		log.Fatalln("Couldn't prepare:", err)
	}
	runner.Deploy(tests) // Delete containers when it's the test runner
}
