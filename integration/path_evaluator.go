package integration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func EvaluatePath(pta string) (dir string, filename string, path string, _ error) {
	pta = strings.TrimSpace(pta)

	if isValidPathFile(pta) {
		// working filepath with filename
		dir, filename = filepath.Split(pta)
		return dir, filename, filepath.Join(dir, filename), nil
	} else if s, err := os.Stat(pta); err == nil && s.IsDir() {
		// working filepath whithout filename
		if lastDir := filepath.Base(pta); lastDir != "." && lastDir != "/" && lastDir != "\\" {
			dir = pta
			filename = lastDir + ".go" // TODO change if other fileextentions are allowed
			if isValidPathFile(filepath.Join(dir, filename)) {
				return dir, filename, filepath.Join(dir, filename), nil
			}
		}
		return "", "", "", errors.New("bad path")
	} else {
		return "", "", "", errors.New("bad path")
	}
}

func EvaluateNewPath(pta string) (dir string, filename string, path string, _ error) {
	pta = strings.TrimSpace(pta)

	// Split into path and

	// check if path is a file
	if strings.HasSuffix(pta, ".go") {

		// check if file already exists
		if inf, err := os.Stat(pta); err == nil && !inf.IsDir() {
			return "", "", "", errors.New("file already exists")
		}

		// check if folder exists
		dir, filename := filepath.Split(pta)
		if _, err := os.Stat(dir); err != nil {

			// check if subfolder exists
			if inf, err := os.Stat(filepath.Dir(dir)); err == nil && inf.IsDir() {

				// subfolder exists create dir ontop
				if err = os.Mkdir(filepath.Base(dir), 0755); err != nil {
					return "", "", "", fmt.Errorf("failed to create folder: %w", err)
				}
				return dir, filename, pta, nil
			} else {
				return "", "", "", errors.New("one or more subfolders don't exist, or the path is wrong")
			}
		} else {
			// folder exists
			return dir, filename, pta, nil
		}
	} else {
		sdir, dir := filepath.Split(pta)

		// check if file already exists
		if inf, err := os.Stat(filepath.Join(pta, dir+".go")); err == nil && !inf.IsDir() {
			return "", "", "", errors.New("file already exists")
		}

		// check if folder exists
		if _, err := os.Stat(pta); err != nil { // this/folder

			// check if subfolder exists
			if inf, err := os.Stat(sdir); err == nil && inf.IsDir() {

				// subfolder exists create dir ontop
				if err = os.Mkdir(pta, 0755); err != nil {
					return "", "", "", fmt.Errorf("failed to create folder: %w", err)
				}
				return pta, dir + ".go", filepath.Join(pta, dir+".go"), nil
			} else {
				return "", "", "", errors.New("one or more subfolders don't exist, or the path is wrong")
			}
		} else {
			// folder exists
			return pta, dir + ".go", filepath.Join(pta, dir+".go"), nil
		}
	}
}

func isValidPathFile(path string) bool {
	if s, err := os.Stat(path); err == nil && !s.IsDir() {
		return true
	}
	return false
}
