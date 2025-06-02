package mrunner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Run a config using the runner (returns the path the go module was generated at)
func (f Factory) GenerateConfigModule(config string, profile string, deleteContainers bool, deployContainers bool, printFunc func(string)) (module string, cachePath string, err error) {

	// Generate the folder for the config
	configModPath := f.ConfigCacheDirectory(config, profile)
	if err := os.RemoveAll(configModPath); err != nil {
		return "", "", fmt.Errorf("couldn't clear config cache dir: %s", err)
	}
	if err := os.MkdirAll(configModPath, 0755); err != nil {
		return "", "", fmt.Errorf("couldn't create config cache dir: %s", err)
	}

	// Copy config file and replace package name
	printFunc("Creating config folder...")
	_, err = f.CopyToCacheWithReplacedPackage(f.ConfigFile(config), filepath.Join(configModPath, "config.go"), "main")
	if err != nil {
		return "", "", fmt.Errorf("couldn't copy and replace: %s", err)
	}

	// Create the run file for the config
	printFunc("Generating files...")
	runFile := GenerateRunFile(deployContainers, deleteContainers)
	if err := os.WriteFile(filepath.Join(configModPath, "run.go"), []byte(runFile), 0755); err != nil {
		return "", "", fmt.Errorf("couldn't create run file: %s", err)
	}

	// Initialize the module
	mod, err := f.PrepareFolderInCache(configModPath, printFunc)
	if err != nil {
		return "", "", fmt.Errorf("couldn't prepare cache: %s", err)
	}

	return mod, configModPath, nil
}

// Get the file path of a config file by name
func (f Factory) ConfigFile(config string) string {
	config = strings.TrimSuffix(config, ".go")
	return filepath.Join(f.mDir, config+".go")
}
