package mrunner

import (
	"strings"
	"unicode"
)

type Filter interface {
	ScanLine(line string) (bool, string)
}

// Scan all lines of content using filters.
//
// Returns all of the results by filter.
func ScanLines(content string, filters []Filter) map[Filter][]string {
	results := map[Filter][]string{}

	// Scan all lines in the content using the filters
	for line := range strings.Lines(content) {
		line = removeGoModComments(line)
		line = strings.TrimSpace(line)
		for _, filter := range filters {
			found, result := filter.ScanLine(line)
			if found {
				results[filter] = append(results[filter], result)
			}
		}
	}

	return results
}

// A filter for getting the module name from a go.mod file
var FilterModuleName = moduleNameFilter{}

type moduleNameFilter struct{}

func (f moduleNameFilter) ScanLine(line string) (bool, string) {

	// Make sure module is at the beginning of the line
	if !strings.HasPrefix(line, "module") {
		return false, ""
	}

	// Get the module name by removing all content until the space
	name := strings.TrimLeftFunc(line, func(r rune) bool {
		return !unicode.IsSpace(r)
	})
	return true, strings.TrimSpace(name)
}

// A filter for getting the module name from a go.mod file
var FilterGoVersion = goVersionFilter{}

type goVersionFilter struct{}

func (f goVersionFilter) ScanLine(line string) (bool, string) {

	// Make sure module is at the beginning of the line
	if !strings.HasPrefix(line, "go") {
		return false, ""
	}

	// Get the module name by removing all content until the space
	version := strings.TrimLeftFunc(line, func(r rune) bool {
		return !unicode.IsSpace(r)
	})
	return true, strings.TrimSpace(version)
}

// Helper function for removing // comments in the go.mod file
func removeGoModComments(line string) string {
	count := 0
	trimmed := strings.TrimRightFunc(line, func(r rune) bool {
		if r == '/' {
			count++
			if count == 2 {
				return false
			}
		} else {
			count = 0
		}
		return true
	})
	if count != 2 {
		return line
	}
	return trimmed[:len(trimmed)-1]
}
