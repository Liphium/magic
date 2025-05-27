package mrunner

import (
	"strings"
)

type Replacer interface {
	Replace(old string) string
}

// Replace all results of a filter with something
func ReplaceLines(content string, replacer Replacer) string {
	newContent := ""

	// Replace all lines in the content using the replacers
	for line := range strings.Lines(content) {
		line = strings.Trim(line, "\n")
		replaced := replacer.Replace(line)
		if replaced == "" {
			continue
		}
		newContent += replaced + "\n"
	}

	return newContent
}

// Empty replacer that doesn't do anything. Useful for streamlining code.
type NoReplacer struct{}

func (r NoReplacer) Replace(old string) string {
	return old
}
