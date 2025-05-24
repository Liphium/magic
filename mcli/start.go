package main

import (
	"errors"
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
	mrunner.GenConfig(path, config, profile)

	// See if the magic directory already exists
	_, err = integration.GetMagicDirectory(3)
	if err != nil {
		return err
	}

	tui.Console.AddItem("Starting...")

	integration.ExecCmdWithFunc(func(s string) {
		tui.Console.AddItem(s)
	}, true,"go", "run", ".")

	tui.RunTui()

	return nil
}
