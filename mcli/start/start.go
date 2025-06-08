package start_command

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Liphium/magic/integration"
	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/magic/mrunner"
	"github.com/Liphium/magic/msdk"
	"github.com/Liphium/magic/tui"
	"github.com/urfave/cli/v3"
)

func BuildCommand() *cli.Command {
	var startConfig = ""
	var startProfile = ""
	return &cli.Command{
		Name:        "start",
		Description: "Magically start your project.",
		Action: func(ctx context.Context, c *cli.Command) error {
			return startCommand(startConfig, startProfile)
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
		},
	}
}

// Command: magic start
func startCommand(config string, profile string) error {
	wbOld, err := os.Getwd()
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
	logLeaf := tui.NewStringLeaf()
	quitLeaf := tui.NewLeaf[error]()
	commandLeaf := tui.NewStringLeaf()
	exitLeaf := tui.NewLeaf[func()]()

	go func() {
		logLeaf.Println("Starting...")
		err := integration.BuildThenRun(func(s string) {
			if strings.HasPrefix(s, mrunner.PlanPrefix) {
				mconfig.CurrentPlan, err = mconfig.FromPrintable(strings.TrimLeft(s, mrunner.PlanPrefix))
				if err != nil {
					logLeaf.Println(strings.TrimLeft(s, mrunner.PlanPrefix))
					quitLeaf.Append(fmt.Errorf("ERROR: couldn't parse plan: %w", err))
					return
				}
				return
			}
			if strings.HasPrefix(s, msdk.StartSignal) {
				return
			}
			logLeaf.Println(strings.TrimRight(s, "\n"))
		}, func(cmd *exec.Cmd) {
			if err = os.Chdir(wbOld); err != nil {
				quitLeaf.Append(fmt.Errorf("ERROR: couldn't change working directory: %w", err))
			}

			exitLeaf.Append(func() {
				// Create a runner and stop all the containers
				runner, err := mrunner.NewRunnerFromPlan(mconfig.CurrentPlan)
				if err == nil {
					runner.StopContainers()
				}

				if err := cmd.Process.Kill(); err != nil {

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
		}, genDir, mod, config, profile, mDir)
		if err != nil {
			quitLeaf.Append(fmt.Errorf("ERROR: failed to start config: %w", err))
		} else {

			if os.Getenv("MAGIC_NO_END") == "true" {
				return
			}
			quitLeaf.Append(fmt.Errorf("Application finished."))
		}
	}()

	// Config for tui
	greenTeaConfig := &tui.GreenTeaConfig{
		RefreshDelay: 100,
		Commands:     getCommands(logLeaf, quitLeaf, exitLeaf),
		LogLeaf:      logLeaf,
		QuitLeaf:     quitLeaf,
		CommandLeaf:  commandLeaf,
		ExitLeaf:     exitLeaf,
	}

	// Start tui
	tui.StartTui(greenTeaConfig)

	return nil
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

func getCommands(logLeaf *tui.StringLeaf, quitLeaf *tui.Leaf[error], exitLeaf *tui.Leaf[func()]) []*cli.Command {

	// Implement commands
	var testPath string
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
		{
			Name:    "test",
			Usage:   "",
			Aliases: []string{"t"},
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name:        "path",
					Destination: &testPath,
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				if testPath != "" {
					go tui.TestCommand(testPath, logLeaf)
				} else {
					tui.CommandError = "usage: test [path]"
				}
				return nil
			},
		},
	}
	return commands
}
