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

func TestFunctionParameterLookup(t *testing.T) {
	sampleCode := `
package main

func foo(a string, b int) {}
func bar(x int, y int) {}
func baz(a string, b int) {}
func qux(a string, b string) {}
func noParams() {}
`

	cases := []struct {
		params   []string
		expected []string
	}{
		{[]string{"string", "int"}, []string{"foo", "baz"}},
		{[]string{"int", "int"}, []string{"bar"}},
		{[]string{"string", "string"}, []string{"qux"}},
		{[]string{}, []string{"noParams"}},
		{[]string{"float64"}, []string{}},
	}

	for _, c := range cases {
		filter := &mrunner.FilterGoFileFunctionParameter{Parameters: c.params}
		found := mrunner.ScanLinesSanitize(sampleCode, []mrunner.Filter{filter}, &mrunner.CommentCleaner{})[filter]
		if len(found) != len(c.expected) {
			t.Errorf("For params %v, expected %d matches, got %d", c.params, len(c.expected), len(found))
			continue
		}
		for i, name := range c.expected {
			if found[i] != name {
				t.Errorf("For params %v, expected function '%s', got '%s'", c.params, name, found[i])
			}
		}
	}
}
