package mrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	_, err := f.CopyToCacheWithReplacedPackage(f.ConfigFile(config), filepath.Join(configModPath, "config.go"), "main")
	if err != nil {
		return "", fmt.Errorf("couldn't copy and replace: %s", err)
	}

	// Create the run file for the config
	printFunc("Generating files...")
	runFile := GenerateRunFile(deployContainers)
	if err := os.WriteFile(filepath.Join(configModPath, "run.go"), []byte(runFile), 0755); err != nil {
		return "", fmt.Errorf("couldn't create run file: %s", err)
	}

	// Initialize the module
	if err := f.PrepareFolderInCache(configModPath, printFunc); err != nil {
		return "", fmt.Errorf("couldn't prepare cache: %s", err)
	}

	return configModPath, nil
}

// Get the file path of a config file by name
func (f Factory) ConfigFile(config string) string {
	config = strings.TrimSuffix(config, ".go")
	return filepath.Join(f.mDir, config+".go")
}
