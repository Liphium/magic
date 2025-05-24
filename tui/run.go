package tui

import (
)

// Command: run [path]
func runCommand(path string, console *sPipe) error {
	// run script 
	console.AddItem("run script with path: " + path)
	return nil
}
