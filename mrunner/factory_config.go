package mrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Liphium/magic/integration"
)

// Run a config using the runner (returns the path the go module was generated at)
func (f Factory) GenerateConfigModule(config string, profile string, deployContainers bool, printFunc func(string)) (string, error) {

	// Generate the folder for the config
	configModPath := f.ConfigCacheDirectory(config, profile)
	if err := os.MkdirAll(configModPath, 0755); err != nil {
		return "", fmt.Errorf("couldn't create config dir: %s", err)
	}

	// Copy config file and replace package name
	printFunc("Creating config folder...")
	content, err := os.ReadFile(f.ConfigFile(config))
	if err != nil {
		return "", fmt.Errorf("couldn't read config file: %s", err)
	}
	newContent := ReplaceLinesSanitized(string(content), GoPackageReplacer{
		NewPackage: "main",
	}, &CommentCleaner{})
	if err := os.WriteFile(filepath.Join(configModPath, "config.go"), []byte(newContent), 0755); err != nil {
		return "", fmt.Errorf("couldn't write new config file: %s", err)
	}

	// Initialize the module
	printFunc("Initializing module...")
	version, err := f.GenerateModFile(configModPath, printFunc)
	if err != nil {
		return "", fmt.Errorf("couldn't generate go.mod: %s", err)
	}

	// Create the run file for the config
	printFunc("Generating files...")
	runFile := GenerateRunFile(deployContainers)
	runFileName := fmt.Sprintf("run_%s.go", config)
	if err := os.WriteFile(filepath.Join(configModPath, runFileName), []byte(runFile), 0755); err != nil {
		return "", fmt.Errorf("couldn't create run file: %s", err)
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
	if err := os.Chdir(configModPath); err != nil {
		return "", fmt.Errorf("couldn't change working directory to mod: %s", err)
	}
	defer os.Chdir(oldWd) // Change back in case of return (to prevent errors)

	// Add the current module to the go.work
	if err := integration.ExecCmdWithFunc(printFunc, "go", "work", "use", "."); err != nil {
		return "", fmt.Errorf("couldn't add mod to work: %s", err)
	}

	// Import all the dependencies from the go.mod
	if err := integration.ExecCmdWithFunc(printFunc, "go", "mod", "tidy"); err != nil {
		return "", fmt.Errorf("couldn't tidy go.mod: %s", err)
	}
	printFunc("Imported dependencies.")

	return configModPath, nil
}

// Get the file path of a config file by name
func (f Factory) ConfigFile(config string) string {
	config = strings.TrimSuffix(config, ".go")
	return filepath.Join(f.mDir, config+".go")
}
