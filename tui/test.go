package tui

import (
	"path/filepath"

	"github.com/Liphium/magic/integration"
	"github.com/tiemingo/greentea"
)

// Command: test [path]
func TestCommand(fp string, console *greentea.StringLeaf) error {

	// set tests as dir
	mDir, err := integration.GetMagicDirectory(5) // beacause cwd is inside ./magic/cache/config_default
	if err != nil {
		console.Printlnf("failed to get magic dir: %s", err)
		return nil
	}
	fp = filepath.Join(mDir, "tests", fp)

	// verify filepath
	_, filename, _, err := integration.EvaluatePath(fp)
	if err != nil {
		console.Printlnf("can't find %s: %s", fp, err.Error())
		return nil
	}

	// run test
	console.Printlnf("Starting test: %s", filename)
	return nil
}
