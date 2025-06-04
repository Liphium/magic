package init_command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Liphium/magic/integration"
	"github.com/Liphium/magic/mrunner"
)

const defaultTestBase = `package %s

import (
	"fmt"
	"testing"

	"github.com/Liphium/magic/mconfig"
)

func Test%s(t *testing.T, p *mconfig.Plan) {
	fmt.Println("Hello, I'm the greatest wizzard of all time!")
}
`

// Command: magic init test <name>
func initTestCommand(fp string) error {
	if strings.TrimSpace(fp) == "" {
		log.Fatalln("Please specify a path like script1.go or script1/hello.go.")
	}

	// Get magic dir
	mDir, err := integration.GetMagicDirectory(5)
	if err != nil {
		return err
	}

	// Create magic/tests if it doesn't exist
	log.Println("Creating tests folder..")
	testDir := filepath.Dir(filepath.Join(mDir, "tests", fp))
	if err := os.MkdirAll(testDir, 0755); err != nil {
		log.Fatalln("Failed to create scripts folder:", err)
	}

	// Check the file for a module name
	files, err := os.ReadDir(testDir)
	if err != nil {
		log.Fatalln("Failed to read directory:", err)
	}
	packageName := fmt.Sprintf("magic_%s", filepath.Base(testDir))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Read the file
		content, err := os.ReadFile(filepath.Join(testDir, file.Name()))
		if err != nil {
			log.Fatalf("Couldn't read file %q: %s \n", file.Name(), err)
		}

		// Find the package name
		results := mrunner.ScanLinesSanitize(string(content), []mrunner.Filter{mrunner.FilterGoFilePackageName}, &mrunner.CommentCleaner{})
		res, ok := results[mrunner.FilterGoFilePackageName]
		if ok && len(res) == 1 {
			packageName = res[0]
		}
	}

	// Evaluate the filepath
	_, filename, path, err := integration.EvaluateNewPath(filepath.Join(mDir, "tests", fp))
	if err != nil {
		return fmt.Errorf("bad path %q: %w", fp, err)
	}

	// Generate test base
	testName := integration.SnakeToCamelCase(strings.TrimRight(filename, ".go"), true)
	testBase := fmt.Sprintf(defaultTestBase, packageName, testName)

	// Create all files needed
	log.Println("Creating test..")
	if err := integration.CreateFileWithContent(path, testBase); err != nil {
		return err
	}

	// Print success message
	log.Println()
	log.Println("Successfully created test template.")

	return nil
}
