package magic

import (
	"errors"
	"fmt"
	"os"

	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/magic/mrunner"
	"github.com/Liphium/magic/scripting"
	"github.com/Liphium/magic/util"
	"github.com/spf13/pflag"
)

// Run a script using the current factory and context
func (f Factory) runScript(script scripting.Script, ctx *mconfig.Context) error {

	// Make sure the plan file exists
	_, err := os.Stat(f.PlanFile(ctx.Profile()))
	if os.IsNotExist(err) {
		return errors.New("scripts can't be run until the app was started at least once for the profile")
	}

	// Check if there is a magic instance running for the current profile
	if !f.IsProfileLocked(ctx.Profile()) {
		util.Log.Println("WARNING: It seems like Magic isn't running, this may cause weird circumstances for scripts.")
	}

	// Read the plan and parse
	content, err := os.ReadFile(f.PlanFile(ctx.Profile()))
	if err != nil {
		return fmt.Errorf("couldn't read plan file: %s", err)
	}
	plan, err := mconfig.FromPrintable(string(content))
	if err != nil {
		return fmt.Errorf("couldn't parse plan from plan file: %s", err)
	}

	// Create the runner for the script from what we found in the plan file
	runner, err := mrunner.NewRunnerFromPlan(plan)
	if err != nil {
		return fmt.Errorf("couldn't create runner from plan: %s", err)
	}

	// Load environment variables into current application
	util.Log.Println("Loading environment...")
	for key, value := range runner.Plan().Environment {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("couldn't set environment variable %s: %s", key, err)
		}
	}
	util.Log.Println("Successfully prepared everything!")
	fmt.Println()

	// Collect data for the script and run
	arguments := script.Collector(pflag.Args())
	if err := script.Handler(runner, arguments); err != nil {
		return fmt.Errorf("script exited with error: %s", err)
	}

	return nil
}
