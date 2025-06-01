package mrunner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Get a script's location as an absolute path
func (f Factory) ScriptDirectory(script string) string {
	return filepath.Join(f.mDir, "scripts", script)
}

// Generate a folder for a script
//
// runFileFormat should have one %s in it that will be replaced with the run function of the script.
func (f Factory) GenerateScriptFolder(path string, runFileFormat string, printFunc func(string)) (string, error) {
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("can't execute with absolute paths")
	}
	path = filepath.Clean(path)

	// Create the folder for the script
	modName := ScriptPathToSnakeCase(path)
	scriptDir := f.ScriptCacheDirectory(modName)
	if err := os.MkdirAll(scriptDir, 0755); err != nil {
		return "", fmt.Errorf("couldn't create script directory: %s", err)
	}

	// Normalize the path for the script
	ogScriptDir := f.ScriptDirectory(path)
	data, err := os.Stat(ogScriptDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("couldn't get script info 1: %s", err)
		}

		// Try adding .go and checking for it then
		data, err = os.Stat(f.ScriptDirectory(path + ".go"))
		if err != nil {
			return "", fmt.Errorf("couldn't get script info 2: %s", err)
		}

		// Add .go since it succeeded
		path = path + ".go"
	}

	// Scan the directory for functions in case it's a script
	fileToCopy := ""
	functionToCall := ""
	if data.IsDir() {
		files, err := os.ReadDir(ogScriptDir)
		if err != nil {
			return "", fmt.Errorf("couldn't list script directory: %s", err)
		}

		// Go through all of the files to find the function to call
		for _, f := range files {
			if f.IsDir() {
				continue
			}

			// Scan the file for functions taking in a plan
			fn, err := scanScriptFileForFunction(filepath.Join(ogScriptDir, f.Name()))
			if err != nil {
				return "", err
			}

			// Make sure there is only one run function
			if functionToCall != "" {
				return "", errors.New("there can only be one run function in a script")
			}
			functionToCall = fn
			fileToCopy = filepath.Join(ogScriptDir, f.Name())
		}
	} else {
		// Scan the file for functions taking in a plan
		functionToCall, err = scanScriptFileForFunction(f.ScriptDirectory(path))
		if err != nil {
			return "", err
		}
		fileToCopy = f.ScriptDirectory(path)
	}

	// Copy over the script file
	_, err = f.CopyToCacheWithReplacedPackage(fileToCopy, filepath.Join(scriptDir, "script.go"), "main")
	if err != nil {
		return "", fmt.Errorf("couldn't copy and replace: %s", err)
	}

	// Generate the run file
	runFile := fmt.Sprintf(runFileFormat, functionToCall)
	if err := os.WriteFile(filepath.Join(scriptDir, "run.go"), []byte(runFile), 0755); err != nil {
		return "", fmt.Errorf("couldn't create run file: %s", err)
	}

	// Prepare the folder
	if err := f.PrepareFolderInCache(scriptDir, printFunc); err != nil {
		return "", fmt.Errorf("couldn't prepare cache: %s", err)
	}

	return scriptDir, nil
}

// Helper function to scan a script for the function that takes in the plan.
func scanScriptFileForFunction(file string) (string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("couldn't read file %s: %s", file, err)
	}
	filter := &FilterGoFileFunctionParameter{
		Parameters: []string{"*mconfig.Plan"},
	}
	results := ScanLinesSanitize(string(content), []Filter{filter}, &CommentCleaner{})
	if res, ok := results[filter]; ok {

		// Make sure there is only one function
		if len(res) != 1 {
			return "", errors.New("found more than one function")
		}

		return res[0], nil
	}
	return "", nil
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
