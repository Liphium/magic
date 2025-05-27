package mrunner_test

import (
	"testing"

	"github.com/Liphium/magic/mrunner"
)

func TestImportLookups(t *testing.T) {
	sampleCode := `
package main

import "fmt"
import (
	"os"
	myAlias "path/filepath"
	"strings"
)

func main() {
	fmt.Println("Hello")
}
`
	filter := &mrunner.FilterGoFileImports{}
	expectedImports := []string{"fmt", "os", "path/filepath", "strings"}

	foundImports := mrunner.ScanLinesSanitize(sampleCode, []mrunner.Filter{filter}, &mrunner.CommentCleaner{})[filter]

	if len(foundImports) != len(expectedImports) {
		t.Fatalf("Expected %d imports, but found %d", len(expectedImports), len(foundImports))
	}

	for i, expected := range expectedImports {
		if i >= len(foundImports) || foundImports[i] != expected {
			t.Errorf("Expected import '%s', but found '%s'", expected, foundImports[i])
		}
	}
}
