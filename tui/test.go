package tui

import (
	"path/filepath"

	"github.com/Liphium/magic/integration"
)

// Command: test [path]
func testCommand(fp string, console *sPipe) error {

	// set tests as dir
	mDir, err := integration.GetMagicDirectory(5) // beacause cwd is inside ./magic/cache/config_default
	if err != nil {
		console.AddItem(err.Error())
		return nil
	}
	fp = filepath.Join(mDir, "tests", fp)

	// verify filepath
	_, filename, _, err := integration.EvaluatePath(fp)
	if err != nil {
		console.AddItem("can't find " + fp + ": " + err.Error())
		return nil
	}

	// run test
	console.AddItem("Starting test: " + filename)
	return nil
}
