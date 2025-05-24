package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Liphium/magic/integration"
	"github.com/Liphium/magic/mrunner"
	"github.com/Liphium/magic/tui"
)

// Command: magic start
func startCommand(config string, profile string) error {
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
		fmt.Println(filepath.Join(mDir, config+".go"))
		return err
	}
	// generate the cache
	wd, err := mrunner.GenConfig(path, config, profile, func(s string) {
		tui.Console.AddItem(s)
	})
	if err != nil {
		return err
	}

	if err = os.Chdir(wd); err != nil {
		return err
	}

	tui.Console.AddItem("Starting...")

	integration.ExecCmdWithFunc(func(s string) {
		tui.Console.AddItem(s)
	}, true, "go", "run", ".")

	tui.RunTui()

	return nil
}
