package test_command

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	mDir, err := integration.GetMagicDirectory(3)
	if err != nil {
		return err
	}
	log.Println("Preparing...")

	// Create a new factory for getting the test directory
	factory := mrunner.NewFactory(mDir)

	// Find all relevant paths (or don't if not desired)
	paths := []string{path}
	if path == "" {
		paths, err = discoverTestDirectories(factory.TestDirectory("."))
	}

	// Convert all paths to relative paths
	relativePaths := make([]string, len(paths))
	for i, path := range paths {
		relativePath, err := filepath.Rel(factory.TestDirectory("."), path)
		if err != nil {
			return fmt.Errorf("Couldn't convert absolute (%s) to relative path: %s", path, err)
		}
		relativePaths[i] = relativePath
	}

	// Start a test runner that goes through all the paths
	if err := startTestRunner(mDir, relativePaths, config, "test"); err != nil {
		return err
	}

	log.Println("Successfully executed test.")
	return nil
}

// Returns all the directories with tests from the given start directory
func discoverTestDirectories(startDir string) ([]string, error) {
	contents, err := os.ReadDir(startDir)
	if err != nil {
		return nil, err
	}

	// Go through all of the contents and search recursively
	paths := []string{}
	found := false
	for _, file := range contents {

		// Check directories recursively
		if file.IsDir() {
			dirPaths, err := discoverTestDirectories(filepath.Join(startDir, file.Name()))
			if err != nil {
				return nil, err
			}
			paths = append(paths, dirPaths...)
			continue
		}

		// Check if it is a go file
		if strings.HasSuffix(file.Name(), ".go") {
			found = true
		}
	}

	// If there was a go file, add the directory
	if found {
		paths = append(paths, startDir)
	}

	return paths, nil
}

// Helper function for starting a new test runner. Can't be run in parallel.
func startTestRunner(mDir string, paths []string, config string, profile string) error {
	// Get the current working directory
	oldWd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create all the folders and stuff
	var mod string
	config, _, mod, err = start_command.CreateStartEnvironment(config, profile, mDir, true)
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

	// Run the tests for each path
	for _, path := range paths {
		loggablePath := path
		if loggablePath == "." || loggablePath == "" {
			loggablePath = "default directory"
		}
		log.Println("Preparing tests in " + loggablePath + "...")

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

		log.Println("Running tests in " + loggablePath + "...")

		// Make a new version of the plan as printable
		printable, err := mconfig.CurrentPlan.ToPrintable()
		if err != nil {
			return fmt.Errorf("couldn't generate printable plan: %s", err)
		}

		// Run go test with the arguments
		if err := integration.ExecCmdWithFunc(func(s string) {
			log.Println(s)
		}, "go", "test", "-args", "plan:"+printable); err != nil {
			return fmt.Errorf("test command failed: %s", err)
		}
	}

	return nil
}
