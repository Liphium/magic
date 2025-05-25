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

// wd: ./magic/cache/
func GoToCache() error {
	err := integration.CreateCache()
	if err != nil {
		return err
	}
	mDir, err := integration.GetMagicDirectory(5)
	if err != nil {
		return err
	}

	err = os.Chdir(filepath.Join(mDir, "cache"))
	if err != nil {
		return fmt.Errorf("failed to cd into %q: %w", filepath.Join(mDir, "cache"), err)
	}
	return nil
}

// wd: ./magic/cache/test_name or script_name if in cache befor
func GenTSFolder(path string, isTest bool) error {
	fName := "script_" + filepath.Base(path)
	if isTest {
		fName = "test_" + filepath.Base(path)
	}

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
		if entry.IsDir() && entry.Name() == fName {
			folderEx = true
		}
	}
	if !folderEx {
		if err := os.Mkdir(fName, 0755); err != nil {
			return err
		}
	}

	wd = filepath.Join(wd, fName)
	err = os.Chdir(wd)
	if err != nil {
		return err
	}
	return nil
}

// wd: ./magic/config/config_profile if in cache folder before
func GenConfFolder(configPath string, confName string) (string, error) {

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
	return wd, nil
}

// no wd change
func CopyFileReplaceModule(fp string, orgName string, newName string) error {
	// COPY configfile and change package
	// Open the file for reading
	file, err := os.Open(fp)
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
	content = strings.Replace(content, "package "+orgName, "package "+newName, 1)

	// Write the replaced content to the file
	err = os.WriteFile(filepath.Join(filepath.Base(fp)), []byte(content), 0755)
	if err != nil {
		return err
	}
	return nil
}

// no wd change
func GenGoMod(isConfig bool, confName string, printFunc func(string)) (string, error) {

	mDir, err := integration.GetMagicDirectory(5)
	if err != nil {
		return "", err
	}

	// load go.mod from conf
	baseDir := filepath.Dir(mDir)

	// Open the file for reading
	file, err := os.Open(filepath.Join(baseDir, "go.mod"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the file content
	moduleName := ""
	version := ""
	scanner := bufio.NewScanner(file)
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

	// init go module
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	integration.ExecCmdWithFunc(printFunc, "go", "mod", "init", filepath.Base(wd))

	// add replace to go.mod
	toadd := "\nreplace " + moduleName + " => ../../../"
	if os.Getenv("MAGIC_DEBUG") == "true" {
		toadd += fmt.Sprintf("\nreplace github.com/Liphium/magic/mconfig => %s", os.Getenv("MAGIC_MCONFIG"))
		toadd += fmt.Sprintf("\nreplace github.com/Liphium/magic/mrunner => %s", os.Getenv("MAGIC_MRUNNER"))
		toadd += fmt.Sprintf("\nreplace github.com/Liphium/magic/integration => %s", os.Getenv("MAGIC_INTEGRATION")) // Add or else
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
	return version, nil
}

func CreateRunFile(deployConainers bool) error {
	// gen runfile
	fc := GenerateRunFile(deployConainers)
	// Open the file for writing (this will truncate the file)
	file, err := os.Create("run.go")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the modified content back to the file
	_, err = file.WriteString(fc)
	if err != nil {
		return err
	}
	return nil
}

func ImportDependencies(confName string, printFunc func(string)) error {
	err := integration.ExecCmdWithFunc(printFunc, "go", "get", confName)
	if err != nil {
		return err
	}
	err = integration.ExecCmdWithFunc(printFunc, "go", "mod", "tidy")
	if err != nil {
		return err
	}
	err = integration.ExecCmdWithFunc(printFunc, "go", "work", "use", ".")
	if err != nil {
		return err
	}
	return nil
}

func GenRunConfig(configPath string, config string, profile string, deployConainers bool, printFunc func(string)) (string, error) {

	confName := config + "_" + profile

	if err := GoToCache(); err != nil {
		return "", err
	}

	printFunc("Creating config folder..")
	if _, err := GenConfFolder(configPath, confName); err != nil {
		return "", err
	}

	if err := CopyFileReplaceModule(configPath, "config", "main"); err != nil {
		return "", err
	}

	printFunc("Initialized module..")
	version, err := GenGoMod(true, confName, printFunc)
	if err != nil {
		return "", err
	}

	printFunc("Creating runfile..")
	if err := CreateRunFile(true); err != nil {
		return "", err
	}

	// Generate the go.work file for this to work
	printFunc("Creating go.work file..")
	if err = handleWorkFile(version); err != nil {
		return "", err
	}

	printFunc("Importing dependencies")
	if err := ImportDependencies(confName, printFunc); err != nil {
		return "", err
	}
	printFunc("Imported dependencies.")

	// Return the directory of the config
	wd, err := os.Getwd()
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
