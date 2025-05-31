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

	return "", nil
}

// Get the file path of a config file by name
func (f Factory) ConfigFile(config string) string {
	config = strings.TrimSuffix(config, ".go")
	return filepath.Join(f.mDir, config+".go")
}
