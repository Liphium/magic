package mrunner

import (
	"strings"
)

type Filter interface {
	Scan(content string) (bool, string)
}

// Scan all lines of content using filters.
//
// Returns all of the results by filter.
func ScanLines(content string, filters []Filter) map[Filter][]string {
	return ScanLinesSanitize(content, filters, NoReplacer{})
}

// Scan all lines of content using filters. First cleans the lines using the sanitizer.
//
// Returns all of the results by filter.
func ScanLinesSanitize(content string, filters []Filter, sanitizer Replacer) map[Filter][]string {
	results := map[Filter][]string{}

	// Fill to make sure the map isn't empty
	for _, filter := range filters {
		results[filter] = []string{}
	}

	// Scan all lines in the content using the filters
	for line := range strings.Lines(content) {
		line = sanitizer.Replace(line)
		if line == "" {
			continue
		}
		for _, filter := range filters {
			found, result := filter.Scan(line)
			if found {
				results[filter] = append(results[filter], result)
			}
		}
	}

	return results
}
