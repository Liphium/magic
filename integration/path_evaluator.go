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

	// check if path is a file
	if !strings.HasSuffix(pta, ".go") {

		// extend path with filename
		if base := filepath.Base(pta); base != "." {
			pta += base + ".go"
		} else {
			return "", "", "", errors.New("")
		}
	}

	if dE, err := DoesDirExist(filepath.Dir(pta)); err != nil {
		return "", "", "", err
	} else if dE {
		return filepath.Dir(pta), filepath.Base(pta), pta, nil
	} else {
		if err = os.MkdirAll(filepath.Dir(pta), 0755); err != nil {
			return "", "", "", fmt.Errorf("failed to create path %q: %w", filepath.Dir(pta), err)
		}
		return filepath.Dir(pta), filepath.Base(pta), pta, nil
	}
}

func isValidPathFile(path string) bool {
	if s, err := os.Stat(path); err == nil && !s.IsDir() {
		return true
	}
	return false
}
