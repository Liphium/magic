package mrunner

import (
	"strings"
	"unicode"
)

// A filter for getting the module name from a go.mod file
var FilterModFileModuleName = moduleNameFilter{}

type moduleNameFilter struct{}

func (f moduleNameFilter) Scan(line string) (bool, string) {

	// Make sure module is at the beginning of the line
	if !strings.HasPrefix(line, "module") {
		return false, ""
	}

	// Remove the module prefix and check if it truely is the module name
	line = strings.TrimPrefix(line, "module")
	if len(line) == 0 || !unicode.IsSpace(rune(line[0])) {
		return false, ""
	}

	// Get the module name by removing all content until the space
	return true, strings.TrimSpace(line)
}

// A filter for getting the module name from a go.mod file
var FilterModFileGoVersion = goVersionFilter{}

type goVersionFilter struct{}

func (f goVersionFilter) Scan(line string) (bool, string) {

	// Make sure module is at the beginning of the line
	if !strings.HasPrefix(line, "go") {
		return false, ""
	}

	// Remove the go prefix to check if it really is the version
	line = strings.TrimPrefix(line, "go")
	if !unicode.IsSpace(rune(line[0])) {
		return false, ""
	}

	// Get the module name by removing all content until the space
	return true, strings.TrimSpace(line)
}

// A filter for getting all replacers and where they point (separated by ;).
//
// Replacers can also contain versions (e.g. github.com/Liphium/chat v1.0.0;github.com/Liphium/magic v1.1.1).
var FilterModFileReplacers = replacerFilter{}

type replacerFilter struct{}

func (f replacerFilter) Scan(line string) (bool, string) {
	line = strings.TrimSpace(line)

	// Check if the line starts with "replace"
	if !strings.HasPrefix(line, "replace") {
		return false, ""
	}

	// Convert to everything expect for the statement itself
	replacer := strings.TrimLeftFunc(line, func(r rune) bool {
		return !unicode.IsSpace(r)
	})

	// Get the origin and destination
	split := strings.Split(replacer, "=>")
	if len(split) != 2 {
		return false, ""
	}
	return true, strings.TrimSpace(split[0]) + ";" + strings.TrimSpace(split[1])
}
