package main

import (
)

// Command: run [path]
func runCommand(path string, console *sPipe) error {
	// run script 
	console.addItem("run script with path: " + path)
	return nil
}
