package main

import (
	"errors"
	"log"
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
	wd, err := mrunner.GenRunConfig(path, config, profile, true, func(s string) {
		log.Println(s)
	})
	if err != nil {
		log.Fatalln("couldn't generate config:", err)
	}
	if err = os.Chdir(wd); err != nil {
		log.Fatalln("couldn't get working directory:", err)
	}

	go func() {
		tui.Console.AddItem("Starting...")
		err := integration.ExecCmdWithFunc(func(s string) {
			tui.Console.AddItem(s)
		}, "go", "run", ".", config, profile)
		if err != nil {
			tui.Console.AddItem("mgc_pan:" + err.Error())
		}
	}()

	tui.RunTui()

	return nil
}
