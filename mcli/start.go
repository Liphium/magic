package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

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
	if config == "" {
		config = "config"
	}
	if profile == "" {
		profile = "default"
	}

	if !integration.IsPathSanitized(config) {
		return errors.New("filename conatins forbidden chars")
	}
	if !integration.IsPathSanitized(profile) {
		return errors.New("profile conatins forbidden chars")
	}

	_, _, path, err := integration.EvaluatePath(filepath.Join(mDir, config+".go"))
	if err != nil {
		return err
	}

	// generate the cache
	wd, err := mrunner.GenRunConfig(path, config, profile, true, func(s string) {
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

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				<-sigs
				cmd.Process.Kill()
				os.Exit(0)
			}()
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
