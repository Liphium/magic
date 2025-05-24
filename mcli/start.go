package main

import (
	"bufio"
	"context"
	"os/exec"

	"github.com/Liphium/magic/tui"
	"github.com/urfave/cli/v3"
)

// Command: magic start
func startCommand(ctx context.Context, c *cli.Command) error {
	tui.Console.AddItem("Starting...")
	cmd := exec.Command("go", "run", ".")

	// Set up the output pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		tui.Console.AddItem(err.Error())
		return nil
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			tui.Console.AddItem(scanner.Text())
		}
	}()

	// Start the command
	if err := cmd.Start(); err != nil {
		tui.Console.AddItem(err.Error())
		return nil
	}

	tui.RunTui()

	return nil
}
