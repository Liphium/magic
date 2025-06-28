package mrunner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Get a test in the original directory
func (f Factory) TestDirectory(path string) string {
	return filepath.Join(f.mDir, "tests", path)
}

// Generate the folder for a test
//
// Returns the directory where the test was generated.
func (f Factory) GenerateTestFolder(path string, printFunc func(string)) (string, error) {
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("can't execute with absolute paths")
	}
	path = filepath.Clean(path)

	// Clean the folder
	modName := ScriptPathToSnakeCase(path)
	testDir := f.TestCacheDirectory(modName)
	if err := os.RemoveAll(testDir); err != nil {
		return "", fmt.Errorf("couldn't clean test directory: %s", err)
	}

	// Create the folder for the test
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return "", fmt.Errorf("couldn't create test directory: %s", err)
	}

	// Normalize the path for the test
	ogTestDir := f.TestDirectory(path)
	data, err := os.Stat(ogTestDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("couldn't get test info 1: %s", err)
		}

		// Try adding .go and checking for it then
		data, err = os.Stat(f.TestDirectory(path + ".go"))
		if err != nil {
			return "", fmt.Errorf("couldn't get test info 2: %s", err)
		}

		// Add .go since it succeeded
		path = path + ".go"
	}

	// Scan the files/file for the test functions
	testMap := map[string][]string{}
	if data.IsDir() {
		files, err := os.ReadDir(ogTestDir)
		if err != nil {
			return "", fmt.Errorf("couldn't list test directory: %s", err)
		}

		// Go through all of the files to find the function to call
		for _, f := range files {
			if f.IsDir() {
				continue
			}

			// Make sure there there are no test files in the directory
			if strings.HasSuffix(f.Name(), "_test.go") {
				// TODO: Evaluate if we actually need this or if it can be removed (might be annoying)
				return "", errors.New("found regular go tests in the test directory")
			}

			// Make sure there is no file named run_test.go (would conflict with the generated test file)
			if f.Name() == "run_test.go" {
				return "", errors.New("found a run_test.go file in your test directory: not allowed due to collision with run file")
			}

			// Scan the file for test functions
			functions, err := scanTestFileForFunctions(filepath.Join(ogTestDir, f.Name()))
			if err != nil {
				return "", err
			}
			if len(functions) == 0 {
				continue
			}

			// Add all of the functions for the file
			testMap[f.Name()] = functions
		}
	} else {
		// Scan the file for test functions
		functions, err := scanTestFileForFunctions(f.TestDirectory(path))
		if err != nil {
			return "", err
		}
		testMap[filepath.Base(path)] = functions
	}

	if data.IsDir() {
		// Copy everything over (in case it's a directory)
		if err := os.CopyFS(testDir, os.DirFS(ogTestDir)); err != nil {
			return "", fmt.Errorf("couldn't copy files to cache: %s", err)
		}

		// Replace the package on all the files in the new directory
		files, err := os.ReadDir(ogTestDir)
		if err != nil {
			return "", fmt.Errorf("couldn't list test directory: %s", err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			startFile := filepath.Join(testDir, file.Name())
			endFile := filepath.Join(testDir, file.Name())
			_, err = f.CopyToCacheWithReplacedPackage(startFile, endFile, "main_test")
			if err != nil {
				return "", fmt.Errorf("couldn't copy and replace: %s", err)
			}
		}

		// Rename all of the files to test files
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			startFile := filepath.Join(testDir, file.Name())
			nameWithoutGo, _ := strings.CutSuffix(file.Name(), ".go")
			endFile := filepath.Join(testDir, nameWithoutGo+"_test.go")
			if err := os.Rename(startFile, endFile); err != nil {
				return "", fmt.Errorf("couldn't rename file %q: %s", file.Name(), err)
			}
		}
	} else {

		// Copy over the script file
		_, err = f.CopyToCacheWithReplacedPackage(f.TestDirectory(path), filepath.Join(testDir, "test_test.go"), "main_test")
		if err != nil {
			return "", fmt.Errorf("couldn't copy and replace: %s", err)
		}
	}

	// Generate the run file
	runFile := testRunFile(testMap)
	if err := os.WriteFile(filepath.Join(testDir, "run_test.go"), []byte(runFile), 0755); err != nil {
		return "", fmt.Errorf("couldn't create run file: %s", err)
	}

	// Prepare the folder
	if _, err := f.PrepareFolderInCache(testDir, printFunc); err != nil {
		return "", fmt.Errorf("couldn't prepare cache: %s", err)
	}

	return testDir, nil
}

// Helper function to scan a test for the functions that take in the plan and testing.T.
func scanTestFileForFunctions(file string) ([]string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file %s: %s", file, err)
	}
	filter := &FilterGoFileFunctionParameter{
		Parameters: []string{"*testing.T", "*mconfig.Plan"},
	}
	results := ScanLinesSanitize(string(content), []Filter{filter}, &CommentCleaner{})
	res, ok := results[filter]
	if !ok {
		return nil, errors.New("something went wrong during scanning of the files")
	}
	return res, nil
}

const defaultTestRunFile = `// Automagically generated by Magic
package main_test

import (
	"os"
	"strings"
	"log"
	"testing"
	
	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/magic/mrunner"
)

func Test(t *testing.T) {

	// Find the plan in the arguments
	printablePlan := ""
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "plan:") {
			printablePlan, _ = strings.CutPrefix(arg, "plan:")
		}
	}

	// Make sure the plan is valid
	if printablePlan == "" {
		log.Fatalln("Couldn't find plan in start arguments.")
	}

	// Parse the plan
	plan, err := mconfig.FromPrintable(printablePlan)
	if err != nil {
		log.Fatalln("Couldn't parse printable plan:", err)
	}

	// Create a new runner from the plan for deleting the databases later
	r, err := mrunner.NewRunnerFromPlan(plan)
	if err != nil {
		log.Fatalln("Couldn't create runner from plan:", err)
	}

	// Run all of the tests
	r.ClearDatabases()
	%s
}
`

const generatedTest = `
	t.Run(%q, func(t *testing.T) {
		%s(t, plan)
	})

	r.ClearDatabases()
`

func testRunFile(testFunctions map[string][]string) string {

	// Generate all of the test functions
	generated := ""
	for file, functions := range testFunctions {
		for _, fn := range functions {
			testName := fmt.Sprintf("%s/%s", file, fn)
			generated += fmt.Sprintf(generatedTest, testName, fn)
		}
	}

	// Generate the final file
	return fmt.Sprintf(defaultTestRunFile, generated)
}
