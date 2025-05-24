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

	err = os.Chdir("cache")
	if err != nil {
		fmt.Println(os.Getwd())
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
	integration.PrintCurrentDirAll()

	// Write the replaced content to the file
	err = os.WriteFile(filepath.Join(wd, config+".go"), []byte(content), 0755)
	if err != nil {
		return "", err
	}

	// load go.mod from conf
	mDir, err := integration.GetMagicDirectory(5)
	if err != nil {
		return "", err
	}
	baseDir := filepath.Dir(mDir)

	// Open the file for reading
	file, err = os.Open(filepath.Join(baseDir, "go.mod"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the file content
	content = ""
	moduleName := ""
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

	// create base go.mod
	integration.ExecCmdWithFunc(printFunc, false, "go", "mod", "init", confName)

	// add replace to go.mod
	toadd := "\nreplace " + moduleName + " => ../../../"

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

	err = integration.ExecCmdWithFunc(printFunc, false, "go", "get", confName)
	if err != nil {
		return "", err
	}
	err = integration.ExecCmdWithFunc(printFunc, false, "go", "get", "githum.com/Liphium/migic/mconfig")
	if err != nil {
		return "", err
	}
	err = integration.ExecCmdWithFunc(printFunc, false, "go", "get", "githum.com/Liphium/migic/mrunner")
	if err != nil {
		return "", err
	}
	err = integration.ExecCmdWithFunc(printFunc, false, "go", "mod", "tidy")
	if err != nil {
		return "", err
	}

	wd, err = os.Getwd()
	if err != nil {
		return "", err
	}

	return wd, nil
}
