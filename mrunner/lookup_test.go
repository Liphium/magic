package mrunner_test

import (
	"fmt"
	"testing"

	"github.com/Liphium/magic/mrunner"
)

const simpleGoMod = `module %s

go %s`

func TestScanLines(t *testing.T) {
	t.Run("module filter", func(t *testing.T) {
		modules := []string{
			"github.com/Liphium/magic/tui",
			"github.com/Liphium/station",
			"github.com/Liphium/magic",
		}

		for i, mod := range modules {
			t.Run(fmt.Sprintf("mod name %d", i), func(t *testing.T) {
				goMod := fmt.Sprintf(simpleGoMod, mod, "1.23")
				results := mrunner.ScanLines(goMod, []mrunner.Filter{mrunner.FilterModuleName})
				if results[mrunner.FilterModuleName][0] != mod {
					t.Fatalf("expected module name %s, got %s", mod, results[mrunner.FilterModuleName][0])
				}
			})
		}
	})

	t.Run("version filter", func(t *testing.T) {
		versions := []string{
			"1.12",
			"2.3.0",
			"5.06",
		}

		for i, ver := range versions {
			t.Run(fmt.Sprintf("version %d", i), func(t *testing.T) {
				goMod := fmt.Sprintf(simpleGoMod, "github.com/Liphium/magic", ver)
				results := mrunner.ScanLines(goMod, []mrunner.Filter{mrunner.FilterGoVersion})
				if results[mrunner.FilterGoVersion][0] != ver {
					t.Fatalf("expected version %s, got %s", ver, results[mrunner.FilterGoVersion][0])
				}
			})
		}
	})

	t.Run("multiple filters", func(t *testing.T) {
		goMod := fmt.Sprintf(simpleGoMod, "github.com/example/test", "1.21")
		filters := []mrunner.Filter{mrunner.FilterModuleName, mrunner.FilterGoVersion}
		results := mrunner.ScanLines(goMod, filters)

		if results[mrunner.FilterModuleName][0] != "github.com/example/test" {
			t.Fatalf("expected module name github.com/example/test, got %s", results[mrunner.FilterModuleName][0])
		}
		if results[mrunner.FilterGoVersion][0] != "1.21" {
			t.Fatalf("expected version 1.21, got %s", results[mrunner.FilterGoVersion][0])
		}
	})

	t.Run("empty input", func(t *testing.T) {
		results := mrunner.ScanLines("", []mrunner.Filter{mrunner.FilterModuleName})
		if len(results[mrunner.FilterModuleName]) != 0 {
			t.Fatalf("expected no results for empty input, got %v", results[mrunner.FilterModuleName])
		}
	})

	t.Run("no filters", func(t *testing.T) {
		goMod := fmt.Sprintf(simpleGoMod, "github.com/example/test", "1.21")
		results := mrunner.ScanLines(goMod, []mrunner.Filter{})
		if len(results) != 0 {
			t.Fatalf("expected no results with no filters, got %v", results)
		}
	})

	t.Run("missing module line", func(t *testing.T) {
		goMod := "go 1.21\n\nrequire (\n\tgithub.com/example/dep v1.0.0\n)"
		results := mrunner.ScanLines(goMod, []mrunner.Filter{mrunner.FilterModuleName})
		if len(results[mrunner.FilterModuleName]) != 0 {
			t.Fatalf("expected no module results when module line missing, got %v", results[mrunner.FilterModuleName])
		}
	})

	t.Run("missing version line", func(t *testing.T) {
		goMod := "module github.com/example/test\n\nrequire (\n\tgithub.com/example/dep v1.0.0\n)"
		results := mrunner.ScanLines(goMod, []mrunner.Filter{mrunner.FilterGoVersion})
		if len(results[mrunner.FilterGoVersion]) != 0 {
			t.Fatalf("expected no version results when go version line missing, got %v", results[mrunner.FilterGoVersion])
		}
	})

	t.Run("complex go.mod with comments", func(t *testing.T) {
		complexGoMod := `// This is a comment
module github.com/example/complex // Another annoying comment

go 1.21.0 // Hello

require (
	github.com/example/dep1 v1.0.0
	github.com/example/dep2 v2.1.0 // indirect
)

replace github.com/example/dep1 => ./local/dep1`

		results := mrunner.ScanLines(complexGoMod, []mrunner.Filter{mrunner.FilterModuleName, mrunner.FilterGoVersion})

		if results[mrunner.FilterModuleName][0] != "github.com/example/complex" {
			t.Fatalf("expected module name github.com/example/complex, got %s", results[mrunner.FilterModuleName][0])
		}
		if results[mrunner.FilterGoVersion][0] != "1.21.0" {
			t.Fatalf("expected version 1.21.0, got %s", results[mrunner.FilterGoVersion][0])
		}
	})

	t.Run("whitespace variations", func(t *testing.T) {
		variations := []string{
			"module  github.com/example/test\ngo  1.21",
			"module\tgithub.com/example/test\ngo\t1.21",
			"  module github.com/example/test  \n  go 1.21  ",
		}

		for i, goMod := range variations {
			t.Run(fmt.Sprintf("variation %d", i), func(t *testing.T) {
				results := mrunner.ScanLines(goMod, []mrunner.Filter{mrunner.FilterModuleName, mrunner.FilterGoVersion})

				if results[mrunner.FilterModuleName][0] != "github.com/example/test" {
					t.Fatalf("expected module name github.com/example/test, got %s", results[mrunner.FilterModuleName][0])
				}
				if results[mrunner.FilterGoVersion][0] != "1.21" {
					t.Fatalf("expected version 1.21, got %s", results[mrunner.FilterGoVersion][0])
				}
			})
		}
	})
}
