package start_command

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Liphium/magic/integration"
	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/magic/mrunner"
	"github.com/Liphium/magic/msdk"
	"github.com/Liphium/magic/tui"
	"github.com/tiemingo/greentea"
	"github.com/urfave/cli/v3"
)

func BuildCommand() *cli.Command {
	var startConfig = ""
	var startProfile = ""
	var startWatch = false
	return &cli.Command{
		Name:        "start",
		Description: "Magically start your project.",
		Action: func(ctx context.Context, c *cli.Command) error {
			return startCommand(startConfig, startProfile, startWatch)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "profile",
				Aliases:     []string{"p"},
				Value:       "",
				Destination: &startProfile,
				Usage:       "The profile that should be used (to run multiple instances of the same config).",
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "",
				Destination: &startConfig,
				Usage:       "The path to the config file that should be used.",
			},
			&cli.BoolFlag{
				Name:        "watch",
				Aliases:     []string{"w"},
				Value:       false,
				Destination: &startWatch,
				Usage:       "Watch for changes and restart the project automatically.",
			},
		},
	}
}

// Command: magic start
func startCommand(config string, profile string, watch bool) error {
	wdOld, err := os.Getwd()
	if err != nil {
		return err
	}
	mDir, err := integration.GetMagicDirectory(3)
	if err != nil {
		return err
	}

	// Create all the folders and stuff
	var mod, genDir string
	config, profile, mod, genDir, err = CreateStartEnvironment(config, profile, mDir, false)
	if err != nil {
		return err
	}

	// Configure tui
	logLeaf := greentea.NewStringLeaf()
	quitLeaf := greentea.NewLeaf[error]()
	exitLeaf := greentea.NewLeaf[func()]()
	commandError := &greentea.CommandError{
		CommandError: "",
	}

	go func() {
		processChan := make(chan *exec.Cmd)

		// Append a closing function here to make sure the process is stopped and all the containers are stopped
		exitLeaf.Append(func() {

			if mconfig.CurrentPlan != nil {
				// Create a runner and stop all the containers
				runner, err := mrunner.NewRunnerFromPlan(mconfig.CurrentPlan)
				if err == nil {
					runner.StopContainers()
				}
			}

			// Stop the process in case there
			process, ok := <-processChan
			if !ok {
				return // No process has been started yet, nothing to kill
			}
			if err := process.Process.Kill(); err != nil {

				// test for err process already finished
				if os.ErrProcessDone != err {
					logLeaf.Println("shutdown err:", err)
				} else {
					logLeaf.Println("process already finished")
				}
			} else {
				logLeaf.Println("successfully killed")
			}
		})

		// Create a start function to re-use it for watch mode
		start := func() {
			err := startBuildAndRun(genDir, wdOld, logLeaf, quitLeaf, processChan, mod, config, profile, mDir)
			if err != nil {

				// If we are in watch mode, only print the error to the command line
				if watch {
					if !strings.Contains(err.Error(), "exit status") {
						logLeaf.Println("ERROR: failed to start config:", err)
					}
				} else {
					quitLeaf.Append(fmt.Errorf("ERROR: failed to start config: %w", err))
				}
			} else {
				// Don't end the process when we're watching for changes (it needs to be executed again)
				if watch || os.Getenv("MAGIC_NO_END") == "true" {
					return
				}
				quitLeaf.Append(fmt.Errorf("application finished"))
			}
		}

		if watch {
			logLeaf.Println("Preparing watching...")

			modDir, err := mrunner.NewFactory(mDir).ModuleDirectory()
			if err != nil {
				quitLeaf.Append(fmt.Errorf("couldn't get module directory: %w", err))
				return
			}

			// Create a listener for watching
			listener := integration.HandleWatching(integration.WatchContext[*exec.Cmd, struct{}]{
				Print: func(s string) {
					logLeaf.Println(s)
				},
				Error: quitLeaf.Append,
				Start: func(a struct{}, job **exec.Cmd, c chan *exec.Cmd) error {
					start()
					return nil
				},
				Stop: func(c *exec.Cmd) error {
					return c.Process.Kill()
				},
				RetrievalChannel: processChan,
			}, struct{}{})

			// Start watching
			if err := integration.WatchDirectory(modDir, func() {
				listener(struct{}{}, "Changes detected, rebuilding...")
			}, mDir); err != nil {
				quitLeaf.Append(fmt.Errorf("couldn't watch: %w", err))
			}
		}

		logLeaf.Println("Starting...")
		start()
	}()

	// Config for tui
	greenTeaConfig := &greentea.GreenTeaConfig{
		RefreshDelay: 100,
		Commands:     getCommands(logLeaf, quitLeaf, exitLeaf, commandError),
		LogLeaf:      logLeaf,
		QuitLeaf:     quitLeaf,
		ExitLeaf:     exitLeaf,
		History: &greentea.History{
			Persistent:    true,
			SavePath:      filepath.Join(mDir, "cache"),
			HistoryLength: 25,
		},
		CommandError: commandError,
	}

	// Start tui
	greentea.RunTui(greenTeaConfig)

	return nil
}

// Start the build and run the program.
func startBuildAndRun(directory string, wdOld string, logLeaf *greentea.StringLeaf, quitLeaf *greentea.Leaf[error], processChan chan *exec.Cmd, arguments ...string) error {
	return integration.BuildThenRun(integration.RunConfig{

		// Function for printing the stuff returned by the process to the current tui
		Print: func(s string) {

			// If it's a plan, make sure to set the current plan from it
			if strings.HasPrefix(s, mrunner.PlanPrefix) {
				var err error
				mconfig.CurrentPlan, err = mconfig.FromPrintable(strings.TrimLeft(s, mrunner.PlanPrefix))
				if err != nil {
					logLeaf.Println(strings.TrimLeft(s, mrunner.PlanPrefix))
					quitLeaf.Append(fmt.Errorf("ERROR: couldn't parse plan: %w", err))
					return
				}
				return
			}

			// If it's the start signal, don't print it
			if strings.HasPrefix(s, msdk.StartSignal) {
				return
			}

			// Otherwise print to output
			logLeaf.Println(strings.TrimRight(s, "\n"))
		},

		// Make sure to properly kill the process when the tui is closed
		Start: func(cmd *exec.Cmd) {
			if err := os.Chdir(wdOld); err != nil {
				quitLeaf.Append(fmt.Errorf("ERROR: couldn't change working directory: %w", err))
			}

			// In case we want to give the process to the next person, do that
			if processChan != nil {
				processChan <- cmd
			}
		},

		Directory: directory,
		Arguments: arguments,
	})
}

// Create the environment for starting from config and profile arguments
//
// Also changes working directory to the folder generated.
func CreateStartEnvironment(config string, profile string, mDir string, deleteContainers bool) (newConfig string, newProfile string, modName string, directory string, err error) {
	// Make sure config and profile are valid and don't contain weird characters or letters
	if config == "" {
		config = "config"
	}
	if profile == "" {
		profile = "default"
	}
	if !integration.IsPathSanitized(config) {
		return "", "", "", "", errors.New("config path contains forbidden chars")
	}
	if !integration.IsPathSanitized(profile) {
		return "", "", "", "", errors.New("profile contains forbidden chars")
	}

	// Create a new factory for creating the directory
	factory := mrunner.NewFactory(mDir)

	// Generate the folder for running in the cache directory
	mod, wd, err := factory.GenerateConfigModule(config, profile, true, deleteContainers, func(s string) {
		if !strings.HasPrefix(s, "ERROR") && !mconfig.VerboseLogging {
			return
		}

		log.Println(s)
	})
	if err != nil {
		return "", "", "", "", fmt.Errorf("couldn't generate config: %s", err)
	}
	if err = os.Chdir(wd); err != nil {
		return "", "", "", "", fmt.Errorf("couldn't change working directory: %s", err)
	}

	return config, profile, mod, wd, err
}

func getCommands(logLeaf *greentea.StringLeaf, quitLeaf *greentea.Leaf[error], exitLeaf *greentea.Leaf[func()], commandError *greentea.CommandError) []*cli.Command {

	// Implement commands
	commands := []*cli.Command{
		{
			Name:    "run",
			Usage:   "",
			Aliases: []string{"r"},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				go tui.RunCommand(cmd, logLeaf, quitLeaf)
				return nil
			},
		},
	}
	return commands
}
