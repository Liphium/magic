package mrunner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Liphium/magic/integration"
)

func GenConfig(configPath string, config string, profile string) error {

	err := integration.CreateCache()
	if err != nil {
		return err
	}
	err = os.Chdir("cache")
	if err != nil {
		return err
	}

	confName := config + "_" + profile

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	files, err := os.ReadDir(wd)
	if err != nil {
		return err
	}

	folderEx := false

	// Check if conf folder already exists
	for _, entry := range files {
		if entry.IsDir() && entry.Name() == confName {
			// folder already exists check for go.mod
			folderEx = true
		}
	}
	if !folderEx {
		if err := os.Mkdir(confName, 0755); err != nil {
			return err
		}
	}

	wd = filepath.Join(wd, confName)
	err = os.Chdir(wd)
	if err != nil {
		return err
	}


	// COPY configfile and change package
	// Open the file for reading
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file content
	var content string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return err
	}

	// Replace the first occurrence of the old word with the new word
	index := strings.Index(content, "package config")
	if index != -1 {
		content = content[:index] + "package main" + content[index+len("package config"):]
	}

	// Open the file for writing (this will truncate the file)
	file, err = os.OpenFile(config + ".go", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the modified content back to the file
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

	gomod := "module "+ confName +"\n"
	fmt.Print(gomod)


	return nil
}
