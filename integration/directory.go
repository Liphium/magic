package integration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Get the magic directory (as an absolute path)
func GetMagicDirectory(amount int) (string, error) {
	if amount <= 0 {
		return "", errors.New("amount can't be 0 or less")
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for i := 0; i < amount; i++ {

		files, err := os.ReadDir(wd)
		if err != nil {
			return "", err
		}

		foundMg := false
		foundGm := false
		// Find the magic folder
		for _, entry := range files {
			if entry.IsDir() && entry.Name() == "magic" {
				foundMg = true
			} else if !entry.IsDir() && entry.Name() == "go.mod" {
				foundGm = true
			}
		}
		if foundMg {
			return filepath.Join(wd, "magic"), nil
		} else if foundGm {
			return "", fmt.Errorf("can't find magic directory, too far back, found go.mod in: %q", wd)
		}
		wd = filepath.Dir(wd)
	}
	return "", errors.New("can't find magic directory")
}

// Check if a directory exists (argument can also just be a file)
func DoesDirExist(dirPath string) (bool, error) {
	_, err := os.Stat(filepath.Dir(dirPath))
	if err != nil {
		return false, fmt.Errorf("path to dir does not exist: %w", err)
	} else {
		s, err := os.Stat(dirPath)
		if err != nil {
			return true, nil
		} else if !s.IsDir() {
			return false, errors.New("path leads to an existing file not a dir")
		} else {
			return false, nil
		}
	}
}

// Print all files in the current directory (useful for debugging)
func PrintCurrentDirAll() {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	files, _ := os.ReadDir(".")

	// Find the magic folder
	for _, entry := range files {
		fmt.Println(entry.Name())
	}
}

// Convert a path from the go.mod file to an absolute path.
//
// For relative paths to be properly parsed you need to be in the correct directory.
func ModulePathToAbsolutePath(path string) string {
	trimmed := strings.TrimSpace(path)
	if strings.HasPrefix(trimmed, "./") || strings.HasPrefix(trimmed, "../") {
		absolute, err := filepath.Abs(path)
		if err != nil {
			return path
		}
		return absolute
	}
	return path
}
