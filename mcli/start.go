package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Liphium/magic/integration"
	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/magic/mrunner"
	"github.com/Liphium/magic/tui"
)

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

	// Make sure config and profile are valid and don't contain weird characters or letters
	if config == "" {
		config = "config"
	}
	if profile == "" {
		profile = "default"
	}
	if !integration.IsPathSanitized(config) {
		return errors.New("config path conatins forbidden chars")
	}
	if !integration.IsPathSanitized(profile) {
		return errors.New("profile conatins forbidden chars")
	}

	// Create a new factory for creating the directory
	factory := mrunner.NewFactory(mDir)

	// Generate the folder for running in the cache directory
	wd, err := factory.GenerateConfigModule(config, profile, true, func(s string) {
		log.Println(s)
	})
	if err != nil {
		log.Fatalln("couldn't generate config:", err)
	}
	if err = os.Chdir(wd); err != nil {
		log.Fatalln("couldn't change working directory:", err)
	}

	go func() {
		tui.Console.AddItem("Starting...")
		err := integration.ExecCmdWithFuncStart(func(s string) {
			if strings.HasPrefix(s, mrunner.PlanPrefix) {
				mconfig.CurrentPlan, err = mconfig.FromPrintable(strings.TrimLeft(s, mrunner.PlanPrefix))
				if err != nil {
					tui.Console.AddItem(strings.TrimLeft(s, mrunner.PlanPrefix))
					tui.Console.AddItem(fmt.Sprintf(tui.MagicPanicPrefix+"ERROR: couldn't parse plan: %s", err))
					return
				}
				return
			}
			tui.Console.AddItem(s)
		}, func(cmd *exec.Cmd) {
			if err = os.Chdir(wbOld); err != nil {
				tui.Console.AddItem(tui.MagicPanicPrefix + "ERROR: couldn't change working directory: " + err.Error())
			}

			tui.ShutdownHook = func() {
				if err := cmd.Process.Kill(); err != nil {
					fmt.Println("shutdown err:", err)
				} else {
					fmt.Println("successfully killed")
				}
			}
		}, "go", "run", ".", config, profile, mDir)
		if err != nil {
			tui.Console.AddItem(tui.MagicPanicPrefix + "" + err.Error())
		} else {

			if os.Getenv("MAGIC_NO_END") == "true" {
				return
			}
			tui.Console.AddItem(tui.MagicPanicPrefix + "Application finished.")
		}
	}()

	tui.RunTui()

	return nil
}
