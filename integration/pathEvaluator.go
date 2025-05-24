package integration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func EvaluatePath(pta string) (dir string, filename string, path string, _ error) {
	pta = strings.Trim(pta, " ")

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

func isValidPathFile(path string) bool {
	if s, err := os.Stat(path); err == nil && !s.IsDir() {
		return true
	}
	return false
}

func TestEval() {
	fmt.Println(EvaluatePath("./Scripts/script1/"))
	fmt.Println(EvaluatePath("./Scripts/script1"))
	fmt.Println(EvaluatePath("./Scripts/script1.go"))
	fmt.Println(EvaluatePath("./"))
	fmt.Println(EvaluatePath("./Scripts/script1/script7"))
	fmt.Println(EvaluatePath("./Scripts/script1/script7/test.go"))
}
