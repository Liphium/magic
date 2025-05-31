package mrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Generate a folder for a script
func (f Factory) GenerateScriptFolder(path string, printFunc func(string)) error {
	if filepath.IsAbs(path) {
		return fmt.Errorf("can't execute with absolute paths")
	}
	path = filepath.Clean(path)

	// Create the folder for the script
	modName := ScriptPathToSnakeCase(path)
	if err := os.MkdirAll(f.ScriptCacheDirectory(modName), 0755); err != nil {
		return fmt.Errorf("couldn't create script directory: %s", err)
	}

	// Get the module name and version
	/*
		modDir, err := f.ModuleDirectory()
		if err != nil {
			return fmt.Errorf("couldn't find module directory: %s", err)
		}
		mod, ver, err := f.ModuleNameAndVersion()
		if err != nil {
			return fmt.Errorf("couldn't get module name: %s", err)
		}
	*/

	// TODO: Generate the rest of the folder

	return nil
}

// Normalize the name of a script path (from dir/script1 to dir/script1.go).
func NormalizeScriptPath(path string) string {
	if !strings.HasSuffix(path, ".go") {
		path = path + ".go"
	}
	return path
}

// Convert a script path to its snake case name (script/dir/script1.go to script_dir_script1).
func ScriptPathToSnakeCase(path string) string {
	path = strings.TrimSuffix(path, ".go")
	path = strings.ToLower(path)
	re := regexp.MustCompile(`[^a-z0-9_/\\ ]`)
	path = re.ReplaceAllString(path, "")
	path = strings.ReplaceAll(path, " ", "_")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "\\", "_")
	return path
}
