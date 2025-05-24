package main

import (
)

// Command: test [path]
func testCommand(path string, console *sPipe) error {
	// run test 
	console.addItem("run test with path: " + path)
	return nil
}
