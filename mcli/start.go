package main

import (
	"errors"
	"os"
	"strings"

	"github.com/Liphium/magic/integration"
	"github.com/Liphium/magic/mrunner"
	"github.com/Liphium/magic/tui"
)

// Command: magic start
func startCommand(filepath string, profile string) error {
	if filepath == "" {
		filepath = "./config.go"
	}
	if profile == "" {
		profile = "default"
	}
	_, filename, path, err := integration.EvaluatePath(filepath)
	if err != nil {
		return err
	}
	if !integration.IsPathSanitized(filename) {
		return errors.New("filename conatins forbidden chars")
	}
	if !integration.IsPathSanitized(profile) {
		return errors.New("profile conatins forbidden chars")
	}
	config := strings.TrimRight(filename, ".go")

	// generate the cache
	wd, err := mrunner.GenConfig(path, config, profile, func(s string) {
		tui.Console.AddItem(s)
	})
	if err != nil{
		return err
	}

	if err = os.Chdir(wd); err != nil{
		return err
	}

	tui.Console.AddItem("Starting...")

	integration.ExecCmdWithFunc(func(s string) {
		tui.Console.AddItem(s)
	}, true, "go", "run", ".")

	tui.RunTui()

	return nil
}
