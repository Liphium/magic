package mrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Liphium/magic/integration"
)

// Get the cache directory
func (f Factory) CacheDirectory() string {
	return filepath.Join(f.mDir, "cache")
}

// Get the directory of a script in the cache directory.
func (f Factory) ScriptCacheDirectory(script string) string {
	return filepath.Join(f.CacheDirectory(), "script_"+script)
}

// Get the directory of a test in the cache directory.
func (f Factory) TestCacheDirectory(script string) string {
	return filepath.Join(f.CacheDirectory(), "test_"+script)
}

// Get the directory of a config in the cache directory.
func (f Factory) ConfigCacheDirectory(config string, profile string) string {
	return filepath.Join(f.CacheDirectory(), config+"_"+profile)
}

const generatedWorkFile = `go %s
`

// Create work file or update in case already there.
func (f Factory) UpdateCacheWorkFileVersion(version string) error {

	// Create work file if it doesn't exist
	workPath := filepath.Join(f.CacheDirectory(), "go.work")
	contentBytes, err := os.ReadFile(workPath)
	if err != nil {

		// Create the work file with the version from the go mod
		if err := os.WriteFile(workPath, []byte(fmt.Sprintf(generatedWorkFile, version)), 0755); err != nil {
			return fmt.Errorf("couldn't write go.work: %s", err)
		}
		return nil
	}

	// Replace the go version in case it already exists
	content := ReplaceLinesSanitized(string(contentBytes), &GoVersionReplacer{
		Version: version,
	}, &OneLineCommentReplacer{})
	if err := os.WriteFile(workPath, []byte(content), 0755); err != nil {
		return fmt.Errorf("couldn't replace go.work: %s", err)
	}

	return nil
}

// Generate the default mod file for any script, test or config.
//
// Returns go version and error.
func (f Factory) GenerateModFile(dir string, printFunc func(string)) (string, string, error) {

	// Get the module directory
	modDir, err := f.ModuleDirectory()
	if err != nil {
		return "", "", fmt.Errorf("couldn't get mod dir: %s", err)
	}

	// Change working directory to module directory
	oldWd, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("couldn't get working directory: %s", err)
	}
	if err := os.Chdir(dir); err != nil {
		return "", "", fmt.Errorf("couldn't change working directory to cache: %s", err)
	}
	defer os.Chdir(oldWd) // Change back to prevent errors

	// Search for module name, go version and replace statement
	content, err := os.ReadFile(filepath.Join(modDir, "go.mod"))
	if err != nil {
		return "", "", fmt.Errorf("couldn't read go.mod: %s", err)
	}
	filters := []Filter{FilterModFileGoVersion, FilterModFileModuleName, FilterModFileReplacers}
	results := ScanLinesSanitize(string(content), filters, &OneLineCommentReplacer{})

	// Make sure the searching actually returned good results
	ver := results[FilterModFileGoVersion]
	mod := results[FilterModFileModuleName]
	replacers := results[FilterModFileReplacers]
	if len(ver) != 1 || len(mod) != 1 {
		return "", "", fmt.Errorf("couldn't find module name or version")
	}

	// Delete the old go mod file (in case there)
	if err := os.RemoveAll(filepath.Join(dir, "go.mod")); err != nil {
		return "", "", fmt.Errorf("couldn't delete go.mod: %s", err)
	}

	// Initialize the new one
	if err := integration.ExecCmdWithFunc(printFunc, "go", "mod", "init", filepath.Base(dir)); err != nil {
		return "", "", fmt.Errorf("couldn't initialize go.mod: %s", err)
	}

	// Change working directory to module directory for the replacer conversion to work properly
	if err := os.Chdir(modDir); err != nil {
		return "", "", fmt.Errorf("couldn't change working directory to mod: %s", err)
	}

	// Put together all the replacers for the new go mod file
	toAdd := ""
	for _, replacer := range replacers {
		args := strings.Split(replacer, ";")

		// Exclude magic debug replacers
		if os.Getenv("MAGIC_DEBUG") == "true" && strings.Contains(args[0], "github.com/Liphium/magic") {
			continue
		}

		// Add the replacer
		toAdd += fmt.Sprintf("\nreplace %s => %s\n", args[0], integration.ModulePathToAbsolutePath(args[1]))
	}

	// Add a replacer for the original module
	toAdd += fmt.Sprintf("\nreplace %s => %s\n", mod[0], modDir)

	// Add additional replacers in debug mode to make sure go doesn't pull from the internet
	if os.Getenv("MAGIC_DEBUG") == "true" {
		toAdd += fmt.Sprintf("\nreplace github.com/Liphium/magic => %s\n", os.Getenv("MAGIC_FOLDER"))
	}

	// Append everything to the new go.mod file
	file, err := os.OpenFile(filepath.Join(dir, "go.mod"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", "", fmt.Errorf("couldn't open go.mod in append mode: %s", err)
	}
	defer file.Close()
	if _, err := file.WriteString(toAdd); err != nil {
		return "", "", fmt.Errorf("couldn't write changes to go.mod: %s", err)
	}

	return mod[0], ver[0], nil
}

// Prepare a new folder in the cache including creating the mod file, etc.
//
// Returns the module name.
func (f Factory) PrepareFolderInCache(directory string, printFunc func(string)) (string, error) {

	// Initialize the module
	printFunc("Initializing module...")
	modName, version, err := f.GenerateModFile(directory, printFunc)
	if err != nil {
		return "", fmt.Errorf("couldn't generate go.mod: %s", err)
	}

	// Update the work file in cache
	if err := f.UpdateCacheWorkFileVersion(version); err != nil {
		return "", fmt.Errorf("couldn't update or generate cache go.work: %s", err)
	}

	// Change working directory to module directory to make sure Go commands don't fail
	printFunc("Importing dependencies...")
	oldWd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("couldn't get working directory: %s", err)
	}
	if err := os.Chdir(directory); err != nil {
		return "", fmt.Errorf("couldn't change working directory to mod: %s", err)
	}
	defer os.Chdir(oldWd) // Change back in case of return (to prevent errors)

	// Add the current module to the go.work
	if err := integration.ExecCmdWithFunc(printFunc, "go", "work", "use", "."); err != nil {
		return "", fmt.Errorf("couldn't add mod to work: %s", err)
	}

	// Add magic as a dependency (only when not in debug mode)
	if os.Getenv("MAGIC_DEBUG") != "true" {
		if err := integration.ExecCmdWithFunc(printFunc, "go", "get", "-u", "github.com/Liphium/magic@"+integration.MagicVersion); err != nil {
			return "", fmt.Errorf("couldn't add Magic as a dependency: %s", err)
		}
	}

	// Import all the dependencies from the go.mod
	if err := integration.ExecCmdWithFunc(printFunc, "go", "mod", "tidy"); err != nil {
		return "", fmt.Errorf("couldn't tidy go.mod: %s", err)
	}
	printFunc("Imported dependencies.")

	return modName, nil
}

// Copy a file to the target directory and replace its package
//
// target should be the full path including the file name
func (f Factory) CopyToCacheWithReplacedPackage(file string, target string, newPackage string) ([]byte, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return content, fmt.Errorf("couldn't read file %s: %s", file, err)
	}
	newContent := ReplaceLinesSanitized(string(content), GoPackageReplacer{
		NewPackage: newPackage,
	}, &CommentCleaner{})
	if err := os.WriteFile(target, []byte(newContent), 0755); err != nil {
		return content, fmt.Errorf("couldn't write new config file: %s", err)
	}
	return content, nil
}
