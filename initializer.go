package magic

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/magic/mrunner"
	"github.com/Liphium/magic/scripting"
	"github.com/Liphium/magic/util"
	"github.com/spf13/pflag"
)

var prepared = false

// Magic flags - defined at package level
var (
	verboseFlag = pflag.Bool("m-verbose", false, "Enable verbose logging for Magic.")
	profileFlag = pflag.String("profile", "default", "The profile to use for Magic.")
	scriptsFlag = pflag.Bool("scripts", false, "Lists all scripts registered in Magic.")
	runFlag     = pflag.StringP("run", "r", "", "Provide this if you want to run a script.")
)

type Config struct {
	// Required. This app name will be used in the container name of databases or other services we may start for you automatically. Make sure you have no other project using the same name.
	AppName string

	// Required. This function will be executed to plan all the containers you want to start. You may create databases and more using the context passed into this function.
	PlanDeployment func(ctx *mconfig.Context)

	// Required. This should start your app like normal. Expect all the database and other containers to be started at this point. Magic will make sure they're all ready by doing health checks and stuff.
	StartFunction func()

	// Add scripts that could be useful while developing this app.
	Scripts []scripting.Script
}

// Start your application with Magic. Make sure to provide all the arguments in the config that are required.
func Start(config Config) {
	// Parse flags before preparing
	pflag.Parse()

	factory := Prepare(config, false)
	if factory == nil {
		fmt.Println()
		return
	}
	fmt.Println()

	// Make sure to unlock the lock file in case the app crashes
	defer func() {
		recover() // Make sure the file is always unlocked, even if the function below panics
		factory.Unlock()
	}()

	// Start the app
	config.StartFunction()
}

// Start all the containers and more. Call this before running tests with Magic.
//
// Returns a factory if execution should continue. Please make sure to unlock the factory when the program exits.
func Prepare(config Config, tests bool) *Factory {

	// Make sure we're not preparing again (this could happen in tests)
	if prepared {
		return nil
	}
	prepared = true

	if config.AppName == "" {
		util.Log.Fatalln("You MUST provide an app name for Magic's config. It will be used in the container name of databases or other services we may start for you automatically. Make sure no two app names are the same as containers might otherwise be deleted without you expecting it.")
	}

	// Enable verbose logging in case desired
	mconfig.VerboseLogging = *verboseFlag

	// Create a new Magic context
	currentProfile := *profileFlag
	if tests {
		currentProfile = "test-" + currentProfile
	}
	ctx := mconfig.DefaultContext(config.AppName, currentProfile)

	// Check if all scripts should be listed
	if *scriptsFlag {
		listScripts(config)
		return nil
	}

	// Create a factory for initializing everything
	factory, err := createFactory()
	if err != nil {
		util.Log.Fatalln("Something went wrong:", err)
		return nil
	}
	factory.WarnIfNotIgnored()

	// Check if a script should be run
	script := *runFlag
	if script != "" {
		i := slices.IndexFunc(config.Scripts, func(s scripting.Script) bool {
			return strings.EqualFold(s.Name, script)
		})
		if i == -1 {
			fmt.Println()
			fmt.Println("This script wasn't found. Here's a list of everything available.")
			listScripts(config)
			return nil
		}

		if err := factory.runScript(config.Scripts[i], ctx); err != nil {
			util.Log.Fatalln("script failure:", err)
		}
		return nil
	}

	// Make sure to lock the profile (to make sure multiple instances of Magic aren't running)
	if err := factory.TryLockProfile(ctx.Profile()); err != nil {
		util.Log.Fatalln("Magic seems to already be running:", err)
	}

	// Plan the deployment
	config.PlanDeployment(ctx)

	// Create the runner and deploy containers
	runner, err := mrunner.NewRunner(ctx)
	if err != nil {
		factory.Unlock()
		util.Log.Fatalln("Couldn't prepare:", err)
	}

	// Write plan to the plan file (we don't want it for the test runner)
	runner.GeneratePlan()
	if !tests {
		plan, err := mconfig.CurrentPlan.ToPrintable()
		if err != nil {
			factory.Unlock()
			util.Log.Fatalln("Failed to create plan:", err)
		}
		if err := os.WriteFile(factory.PlanFile(ctx.Profile()), []byte(plan), 0755); err != nil {
			factory.Unlock()
			util.Log.Fatalln("Failed to write to plan file:", err)
		}
	}

	// Deploy containers (delete containers when it's the test runner)
	if err := runner.Deploy(tests); err != nil {
		factory.Unlock()
		util.Log.Fatalln("Couldn't deploy containers:", err)
	}
	return &factory
}

// Print a list of all the scripts
func listScripts(config Config) {
	fmt.Println()
	fmt.Println("Listing all scripts registered in Magic.")
	fmt.Println("Use -r or --run <script> to run any of them.")
	fmt.Println()
	for _, script := range config.Scripts {
		fmt.Printf("%s - %s \n", script.Name, script.Description)
	}
}
