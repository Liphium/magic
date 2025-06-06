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

func TestFunctionParameterLookupStartsWith(t *testing.T) {
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
		{[]string{"string"}, []string{"foo;string;int", "baz;string;int", "qux;string;string"}},
		{[]string{"int", "int"}, []string{"bar;int;int"}},
		{[]string{"int"}, []string{"bar;int;int"}},
		{[]string{"string", "string"}, []string{"qux;string;string"}},
		{[]string{"float64"}, []string{}},
		{[]string{}, []string{"foo;string;int", "bar;int;int", "baz;string;int", "qux;string;string", "noParams"}},
	}

	for _, c := range cases {
		filter := &mrunner.FilterGoFileFunctionParameter{Parameters: c.params, StartsWith: true}
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

func TestPackageNameFilter(t *testing.T) {
	cases := []struct {
		line     string
		expected bool
		name     string
	}{
		{"package main", true, "main"},
		{"package mrunner", true, "mrunner"},
		{"package 123abc", true, "123abc"},
		{"package", false, ""},
		{"import fmt", false, ""},
		{"func main()", false, ""},
		{"package main extra", false, ""},
	}

	for _, c := range cases {
		ok, name := mrunner.FilterGoFilePackageName.Scan(c.line)
		if ok != c.expected {
			t.Errorf("For line '%s', expected match=%v, got %v", c.line, c.expected, ok)
		}
		if ok && name != c.name {
			t.Errorf("For line '%s', expected name '%s', got '%s'", c.line, c.name, name)
		}
	}
}
