package tui

import (
)

// Command: test [path]
func testCommand(path string, console *sPipe) error {
	// run test 
	console.AddItem("run test with path: " + path)
	return nil
}
