package mrunner

import (
	"strings"
	"unicode"
)

type FilterGoFileImports struct {
	started bool
}

func (f *FilterGoFileImports) Scan(content string) (bool, string) {
	content = strings.TrimSpace(content)

	// If started, scan
	if f.started {

		importCount := 0
		for im := range strings.SplitSeq(content, "\"") {
			importCount++

			// End the searching in case it's closed
			if strings.Contains(im, ")") {
				f.started = false
				return false, ""
			}

			// Return the import in case there is one
			if importCount%2 == 0 {
				return true, im
			}
		}

		return false, ""
	}

	// Check if it starts with import
	if strings.HasPrefix(content, "import") {

		// Remove import from the content
		after := strings.TrimLeftFunc(content, func(r rune) bool {
			return !unicode.IsSpace(r)
		})
		after = strings.TrimSpace(after)

		// Check if a block opens up
		if strings.Contains(after, "(") {
			f.started = true
			return false, ""
		}

		// Scan for import statement
		split := strings.Split(content, "\"")
		if len(split) < 3 {
			return false, ""
		}
		return true, split[1]
	}

	return false, ""
}
