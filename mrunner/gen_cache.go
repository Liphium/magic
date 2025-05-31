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
		return fmt.Errorf("couldn't create cache directory: %s", err)
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

// Expects to be ran with the working directory being the cache folder.
//
// Goes to ./magic/cache/script_name in case successful.
func GenerateScriptFolder(path string) error {
	return generateCacheFolder("script_" + filepath.Base(path))
}

// Expects to be ran with the working directory being the cache folder.
//
// Goes to ./magic/cache/test_name in case successful.
func GenerateTestFolder(path string) error {
	return generateCacheFolder("test_" + filepath.Base(path))
}

// Helper function for generation of folders inside the cache folder.
//
// Expects to be ran in the cache folder.
func generateCacheFolder(name string) error {
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
		if entry.IsDir() && entry.Name() == name {
			folderEx = true
		}
	}
	if !folderEx {
		if err := os.Mkdir(name, 0755); err != nil {
			return err
		}
	}

	wd = filepath.Join(wd, name)
	err = os.Chdir(wd)
	if err != nil {
		return err
	}
	return nil
}

// no wd change
func CopyFileAndReplacePackage(fp string, orgName string, newName string) error {
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
func GenGoMod(mDir string, printFunc func(string)) (string, error) {
	// load go.mod from conf
	baseDir := filepath.Dir(mDir)

	// Open the file for reading
	file, err := os.Open(filepath.Join(baseDir, "go.mod"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the file content
	toAdd := ""
	moduleName := ""
	version := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		if strings.HasPrefix(t, "module ") {
			itms := strings.Split(t, " ")
			ind := slices.Index(itms, "module")
			if ind == -1 || len(itms) < ind+2 {
				return "", errors.New("can't find module name in go.mod")
			}
			moduleName = itms[ind+1]
		}
		if strings.HasPrefix(t, "go ") {
			itms := strings.Split(t, " ")
			ind := slices.Index(itms, "go")
			if ind == -1 || len(itms) < ind+2 {
				return "", errors.New("can't find go version in go.mod")
			}
			version = itms[ind+1]
		}
		if strings.HasPrefix(strings.TrimSpace(t), "replace") {
			if strings.Contains(t, "replace github.com/Liphium/magic/mconfig") || strings.Contains(t, "replace github.com/Liphium/magic/mrunner") || strings.Contains(t, "replace github.com/Liphium/magic/integration") {
				if os.Getenv("MAGIC_DEBUG") == "true" {
					continue
				}
			}
			toAdd += "\n" + strings.Replace(strings.TrimSpace(t), "../", "../../../../", 1)
		}
	}
	if moduleName == "" {
		return "", errors.New("can't find module name in go.mod")
	}

	// add replace to go.mod
	toAdd += "\nreplace " + moduleName + " => ../../../"
	if os.Getenv("MAGIC_DEBUG") == "true" {
		toAdd += fmt.Sprintf("\nreplace github.com/Liphium/magic/mconfig => %s", os.Getenv("MAGIC_MCONFIG"))
		toAdd += fmt.Sprintf("\nreplace github.com/Liphium/magic/mrunner => %s", os.Getenv("MAGIC_MRUNNER"))
		toAdd += fmt.Sprintf("\nreplace github.com/Liphium/magic/integration => %s", os.Getenv("MAGIC_INTEGRATION")) // Add or else
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

	// Open the file in append mode
	file, err = os.OpenFile("go.mod", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write the new line to the file
	if _, err := file.WriteString(toAdd); err != nil {
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

func ImportDependencies(printFunc func(string)) error {
	err := integration.ExecCmdWithFunc(printFunc, "go", "mod", "tidy")
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

	mDir, err := integration.GetMagicDirectory(5)
	if err != nil {
		return "", err
	}

	printFunc("Creating config folder..")
	if err := generateCacheFolder(confName); err != nil {
		return "", err
	}

	if err := CopyFileAndReplacePackage(configPath, "config", "main"); err != nil {
		return "", err
	}

	printFunc("Initialized module..")
	version, err := GenGoMod(mDir, printFunc)
	if err != nil {
		return "", err
	}

	printFunc("Creating runfile..")
	if err := CreateRunFile(true); err != nil {
		return "", err
	}

	// Generate the go.work file for this to work
	printFunc("Creating go.work file..")
	if err = HandleWorkFile(version); err != nil {
		return "", err
	}

	printFunc("Importing dependencies")
	if err := ImportDependencies(printFunc); err != nil {
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

// Create work file or update in case already there.
func HandleWorkFile(version string) error {

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
