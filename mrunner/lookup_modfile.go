package mrunner

import (
	"strings"
	"unicode"
)

// A filter for getting the module name from a go.mod file
var FilterModFileModuleName = moduleNameFilter{}

type moduleNameFilter struct{}

func (f moduleNameFilter) Scan(line string) (bool, string) {
	line = RemoveOneLineComments(line)
	line = strings.TrimSpace(line)

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
var FilterModFileGoVersion = goVersionFilter{}

type goVersionFilter struct{}

func (f goVersionFilter) Scan(line string) (bool, string) {
	line = RemoveOneLineComments(line)
	line = strings.TrimSpace(line)

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
