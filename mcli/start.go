package main

import (
	"errors"
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
		return err
	}
	// generate the cache
	wd, err := mrunner.GenConfig(path, config, profile, true, func(s string) {
		tui.Console.AddItem(s)
	})
	if err != nil {
		return err
	}
	if err = os.Chdir(wd); err != nil {
		return err
	}

	tui.Console.AddItem("Starting...")
	go func() {
		err := integration.ExecCmdWithFunc(func(s string) {
			tui.Console.AddItem(s)
		}, "go", "run", ".", config, profile)
		if err != nil {
			tui.Console.AddItem("mgc_pan:"+err.Error())
		}
	}()

	tui.RunTui()

	return nil
}
