package mrunner

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Liphium/magic/integration"
)

func GenConfig(configPath string, config string, profile string, printFunc func(string)) (string, error) {

	err := integration.CreateCache()
	if err != nil {
		return "", err
	}

	mDir, err := integration.GetMagicDirectory(5)
	if err != nil {
		return "", err
	}

	err = os.Chdir("cache")
	if err != nil {
		return "", err
	}

	confName := config + "_" + profile

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	files, err := os.ReadDir(wd)
	if err != nil {
		return "", err
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
			return "", err
		}
	}

	wd = filepath.Join(wd, confName)
	err = os.Chdir(wd)
	if err != nil {
		return "", err
	}

	// COPY configfile and change package
	// Open the file for reading
	file, err := os.Open(configPath)
	if err != nil {
		return "", err
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
		return "", err
	}

	// Replace the first occurrence of the old word with the new word
	content = strings.Replace(content, "package config", "package main", 1)

	// Write the replaced content to the file
	err = os.WriteFile(filepath.Join(wd, config+".go"), []byte(content), 0755)
	if err != nil {
		return "", err
	}

	// load go.mod from conf
	baseDir := filepath.Dir(mDir)

	// Open the file for reading
	file, err = os.Open(filepath.Join(baseDir, "go.mod"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the file content
	printFunc("Finding version and module name..")
	content = ""
	moduleName := ""
	version := ""
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		if strings.Contains(t, "module ") {
			itms := strings.Split(t, " ")
			ind := slices.Index(itms, "module")
			if ind == -1 || len(itms) < ind+2 {
				return "", errors.New("can't find module name in go.mod")
			}
			moduleName = itms[ind+1]
		}
		if strings.Contains(t, "go ") {
			itms := strings.Split(t, " ")
			ind := slices.Index(itms, "go")
			if ind == -1 || len(itms) < ind+2 {
				return "", errors.New("can't find module name in go.mod")
			}
			version = itms[ind+1]
			break
		}
	}
	if moduleName == "" {
		return "", errors.New("can't find module name in go.mod")
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// create base go.mod and delete old one
	if err := os.RemoveAll("go.mod"); err != nil {
		return "", fmt.Errorf("couldn't delete go.mod: %s", err)
	}
	printFunc("Initializing " + confName + "..")
	integration.ExecCmdWithFunc(printFunc, "go", "mod", "init", confName)

	// add replace to go.mod
	toadd := "\nreplace " + moduleName + " => ../../../"
	if os.Getenv("MAGIC_DEBUG") == "true" {
		toadd += fmt.Sprintf("\nreplace github.com/Liphium/magic/mconfig => %s", os.Getenv("MAGIC_MCONFIG"))
		toadd += fmt.Sprintf("\nreplace github.com/Liphium/magic/mrunner => %s", os.Getenv("MAGIC_MRUNNER"))
	}

	// Open the file in append mode
	file, err = os.OpenFile("go.mod", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write the new line to the file
	if _, err := file.WriteString(toadd); err != nil {
		return "", err
	}

	// gen runfile
	fc := GenerateRunFile(false)
	// Open the file for writing (this will truncate the file)
	file, err = os.Create("run.go")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write the modified content back to the file
	_, err = file.WriteString(fc)
	if err != nil {
		return "", err
	}
	printFunc("Initialized module.")

	// Generate the go.work file for this to work
	printFunc("Creating go.work file..")
	if err := handleWorkFile(version); err != nil {
		return "", fmt.Errorf("couldn't generate work file: %s", err)
	}

	printFunc("Importing dependencies..")
	err = integration.ExecCmdWithFunc(printFunc, "go", "get", confName)
	if err != nil {
		return "", err
	}
	err = integration.ExecCmdWithFunc(printFunc, "go", "mod", "tidy")
	if err != nil {
		return "", err
	}
	err = integration.ExecCmdWithFunc(printFunc, "go", "work", "use", ".")
	if err != nil {
		return "", err
	}
	printFunc("Imported dependencies.")

	// Return the directory of the config
	wd, err = os.Getwd()
	if err != nil {
		return "", err
	}
	return wd, nil
}

const generatedWorkFile = `go %s
`

func handleWorkFile(version string) error {

	// Get the name of the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("couldn't get wd: %s", err)
	}

	// Go to the work file and open
	if err := os.Chdir(".."); err != nil {
		return fmt.Errorf("couldn't go to previous dir: %s", err)
	}
	file, err := os.OpenFile("go.work", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {

		// Create the work file with the version from the go mod
		if err := os.WriteFile("go.work", []byte(fmt.Sprintf(generatedWorkFile, version)), 0755); err != nil {
			return fmt.Errorf("couldn't write go.work: %s", err)
		}
		return nil
	}
	defer file.Close()

	// Modify the version in the file
	content := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(content, "go") {
			content += fmt.Sprintf("go %s\n", version)
		} else {
			content += scanner.Text() + "\n"
		}
	}

	// Go back to the previous directory
	if err := os.Chdir(wd); err != nil {
		return fmt.Errorf("couldn't change back to original wd: %s", err)
	}
	return nil
}
