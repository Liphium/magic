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

const defaultScriptBase = `package %s

import (
	"fmt"

	"github.com/Liphium/magic/mconfig"
)

func Run%s(p *mconfig.Plan) {
	fmt.Println("I'm a wizzard!")
}
`

// Command: magic init script <name>
func initScriptCommand(fp string) error {
	if strings.TrimSpace(fp) == "" {
		log.Fatalln("Please specify a path like script1.go or script1/hello.go.")
	}

	// Get magic dir
	mDir, err := integration.GetMagicDirectory(5)
	if err != nil {
		return err
	}

	// Create magic/scripts if it doesn't exist
	log.Println("Creating scripts folder..")
	scriptDir := filepath.Dir(filepath.Join(mDir, "scripts", fp))
	if err := os.MkdirAll(scriptDir, 0755); err != nil {
		log.Fatalln("Failed to create scripts folder: %w", err)
	}

	// Check the directory for a module name
	files, err := os.ReadDir(scriptDir)
	if err != nil {
		log.Fatalln("Failed to read directory:", err)
	}
	packageName := fmt.Sprintf("magic_%s", filepath.Base(scriptDir))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Read the file
		content, err := os.ReadFile(filepath.Join(scriptDir, file.Name()))
		if err != nil {
			log.Fatalf("Couldn't read file %q: %s \n", file.Name(), err)
		}

		// Find the package name
		results := mrunner.ScanLinesSanitize(string(content), []mrunner.Filter{mrunner.FilterGoFilePackageName}, &mrunner.CommentCleaner{})
		res, ok := results[mrunner.FilterGoFilePackageName]
		if ok && len(res) != 1 {
			packageName = res[0]
		}
	}

	// Evaluate the filepath
	_, filename, path, err := integration.EvaluateNewPath(filepath.Join(mDir, "scripts", fp))
	if err != nil {
		return fmt.Errorf("bad path %s: %w", fp, err)
	}

	// Generate Script base
	scriptName := integration.SnakeToCamelCase(strings.TrimRight(filename, ".go"), true)
	scriptBase := fmt.Sprintf(defaultScriptBase, packageName, scriptName)

	// Create all files needed
	log.Println("Creating script..")
	if err := integration.CreateFileWithContent(path, scriptBase); err != nil {
		return err
	}

	// Print success message
	log.Println()
	log.Println("Successfully created script template.")

	return nil
}
