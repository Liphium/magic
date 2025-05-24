package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/Liphium/magic/integration"
)

// Command: magic init script
func initScriptCommand(fp string) error {

	// Get magic dir
	mDir, err := integration.GetMagicDirectory(5)
	if err != nil {
		return err
	}

	// Evaluate the filepath
	_, filename, path, err := integration.EvaluateNewPath(filepath.Join(mDir, "scripts", fp))
	if err != nil {
		return fmt.Errorf("bad path "+fp+": %w", err)
	}

	// Generate Script base
	var scriptCenter = `
	fmt.Println("I'm a wizzard")
`
	scriptBase := "package scripts\n\nimport(\n    \"fmt\"\n)\n\nfunc run" + strings.TrimRight(filename, ".go") + "(){" + scriptCenter + "}"

	// Create all files needed
	log.Println("Creating script..")
	if err := integration.CreateFileWithContent(path, scriptBase); err != nil {
		return err
	}

	// Print success message
	log.Println()
	log.Println("Successfully created script template.")

	return nil
}
