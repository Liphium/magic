package magic

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/Liphium/magic/v2/mconfig"
	"github.com/Liphium/magic/v2/mrunner"
	"github.com/Liphium/magic/v2/scripting"
	"github.com/Liphium/magic/v2/util"
	"github.com/spf13/pflag"
)

// Magic flags - defined at package level
var (
	verboseFlag = pflag.Bool("m-verbose", false, "Enable verbose logging for Magic.")
	profileFlag = pflag.StringP("profile", "p", "default", "The profile to use for Magic.")
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

	// This is how long Magic waits for your app to finally start up before just killing any test runner. Default: 10 seconds.
	//
	// Hint: You can use magic.Ptr(duration) to convert your duration to a pointer, we use a pointer here to be able to detect it being not set.
	TestAppTimeout *time.Duration
}

// Start all the containers and more. This should be called before doing anything with Magic.
//
// Returns a factory and runner if execution should continue. Please make sure to unlock the factory when the program exits.
func prepare(config Config, testProfile string) (*Factory, *mrunner.Runner) {
	isTestRunner := testProfile != ""

	if config.AppName == "" {
		util.Log.Fatalln("You MUST provide an app name for Magic's config. It will be used in the container name of databases or other services we may start for you automatically. Make sure no two app names are the same as containers might otherwise be deleted without you expecting it.")
	}

	// Enable verbose logging in case desired
	mconfig.VerboseLogging = *verboseFlag

	// Create a new Magic context
	currentProfile := *profileFlag
	if isTestRunner {
		currentProfile = "test-" + testProfile
	}

	// Create a factory for initializing everything
	factory, err := createFactory()
	if err != nil {
		util.Log.Fatalln("Something went wrong:", err)
		return nil, nil
	}
	if mconfig.VerboseLogging {
		util.Log.Println("Using project directory:", factory.projectDir)
	}
	factory.WarnIfNotIgnored()

	// Create the context for Magic config generation
	ctx := mconfig.DefaultContext(config.AppName, currentProfile, factory.projectDir)

	// Check if all scripts should be listed
	if *scriptsFlag && !isTestRunner {
		listScripts(config)
		return nil, nil
	}

	// Check if a script should be run
	script := *runFlag
	if script != "" && !isTestRunner {
		i := slices.IndexFunc(config.Scripts, func(s scripting.Script) bool {
			return strings.EqualFold(s.Name, script)
		})
		if i == -1 {
			fmt.Println()
			fmt.Println("This script wasn't found. Here's a list of everything available.")
			listScripts(config)
			return nil, nil
		}

		if err := factory.runScript(config.Scripts[i], ctx); err != nil {
			util.Log.Fatalln("script failure:", err)
		}
		return nil, nil
	}

	// Make sure to wait until the profile is unlocked (to make sure multiple instances of Magic aren't running)
	var tries int = 0
	var lockErr error = factory.TryLockProfile(ctx.Profile())
	for ; lockErr != nil && errors.Is(lockErr, errProfileLocked); lockErr = factory.TryLockProfile(ctx.Profile()) {
		if tries == 0 {
			fmt.Println()
			util.Log.Println("ERROR: The profile you're trying to run Magic with (" + ctx.Profile() + ") seems to be locked.")
			fmt.Println()
			util.Log.Println("This can mean multiple things:")
			util.Log.Println("1. Most likely: Another instance of Magic is running under the same profile. If you want to run two instances of your app, use the --profile (-p) flag.")
			util.Log.Println("2. Your operating system is not unlocking the lock file, you can try deleting it:", factory.LockFile(ctx.Profile()))
			util.Log.Println("3. Not very likely: A bug in Magic occured that caused this message to appear.")
			fmt.Println()
		}

		if tries%4 == 0 {
			util.Log.Println("Waiting for profile to be unlocked...")
		}

		time.Sleep(500 * time.Millisecond)
		tries++
	}
	if lockErr != nil {
		util.Log.Fatalln("Magic failed to lock the profile:", err)
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
	plan := runner.GeneratePlan()
	if plan != nil && !isTestRunner {
		planStr, err := plan.ToPrintable()
		if err != nil {
			factory.Unlock()
			util.Log.Fatalln("Failed to convert plan:", err)
		}
		if err := os.WriteFile(factory.PlanFile(ctx.Profile()), []byte(planStr), 0755); err != nil {
			factory.Unlock()
			util.Log.Fatalln("Failed to write to plan file:", err)
		}
	}

	// Deploy containers (delete containers when it's the test runner)
	if err := runner.Deploy(isTestRunner); err != nil {
		factory.Unlock()
		util.Log.Fatalln("Couldn't deploy containers:", err)
	}
	return &factory, runner
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

// A little helper to convert anything to a pointer.
func Ptr[T any](obj T) *T {
	return &obj
}
