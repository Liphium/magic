package test_command

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Liphium/magic/integration"
	start_command "github.com/Liphium/magic/mcli/start"
	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/magic/mrunner"
	"github.com/Liphium/magic/msdk"
	"github.com/urfave/cli/v3"
)

func BuildCommand() *cli.Command {
	var testPath = ""
	var startConfig = ""
	return &cli.Command{
		Name:        "test",
		Description: "Magically test your project.",
		Action: func(ctx context.Context, c *cli.Command) error {
			return runTestCommand(testPath, startConfig)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "",
				Destination: &startConfig,
				Usage:       "The path to the config file that should be used.",
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "path",
				Destination: &testPath,
			},
		},
	}
}

// Command: magic test [path]
func runTestCommand(path string, config string) error {
	oldWd, err := os.Getwd()
	if err != nil {
		return err
	}
	mDir, err := integration.GetMagicDirectory(3)
	if err != nil {
		return err
	}
	log.Println("Preparing...")

	// Create all the folders and stuff
	var mod string
	var profile string
	config, profile, mod, err = start_command.CreateStartEnvironment(config, "test", mDir, true)
	if err != nil {
		return err
	}

	// Start the app
	processChan := make(chan *exec.Cmd)
	finishedChan := make(chan struct{})
	go func() {
		if err := integration.ExecCmdWithFuncStart(func(s string) {

			// Wait for a plan to be sent
			if strings.HasPrefix(s, mrunner.PlanPrefix) {
				mconfig.CurrentPlan, err = mconfig.FromPrintable(strings.TrimLeft(s, mrunner.PlanPrefix))
				if err != nil {
					log.Fatalln("Couldn't parse plan:", err)
				}
				return
			}

			// Wait for the start signal from the SDK
			if strings.HasPrefix(s, msdk.StartSignal) {
				finishedChan <- struct{}{}
				return
			}

			// Only print logs when verbose logging
			if !strings.HasPrefix(s, "ERROR") && !mconfig.VerboseLogging {
				return
			}

			log.Println(s)
		}, func(c *exec.Cmd) {
			processChan <- c
		}, "go", "run", ".", mod, config, profile, mDir); err != nil {
			log.Fatalln("couldn't run the app:", err)
		}
	}()

	// Wait for the signal from the SDK to run tests
	<-finishedChan
	process := <-processChan
	defer process.Process.Kill()

	// Go back to the old working directory
	if err := os.Chdir(oldWd); err != nil {
		return fmt.Errorf("couldn't change to old working dir: %s", err)
	}

	// Create a factory for the test creation
	factory := mrunner.NewFactory(mDir)

	// Create the folder for the test
	testCacheDir, err := factory.GenerateTestFolder(path, func(s string) {
		// Only print logs for verbose logging mode
		if !strings.HasPrefix(s, "ERROR") && !mconfig.VerboseLogging {
			return
		}

		log.Println(s)
	})
	if err != nil {
		return fmt.Errorf("couldn't generate test folder: %s", err)
	}

	// Change the working directory to the test folder
	if err := os.Chdir(testCacheDir); err != nil {
		return fmt.Errorf("couldn't go to test dir: %s", err)
	}

	log.Println("Running test " + path + "...")

	// Make a new version of the plan as printable
	printable, err := mconfig.CurrentPlan.ToPrintable()
	if err != nil {
		return fmt.Errorf("couldn't generate printable plan: %s", err)
	}

	// Run go test with the arguments
	if err := integration.ExecCmdWithFunc(func(s string) {
		log.Println(s)
	}, "go", "test", "-args", "plan:"+printable); err != nil {
		return fmt.Errorf("couldn't run test command: %s", err)
	}

	log.Println("Successfully executed test.")
	return nil
}
