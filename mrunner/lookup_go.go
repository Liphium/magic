package mrunner

import (
	"slices"
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

// A filter for searching function names by the parameters the function takes in.
type FilterGoFileFunctionParameter struct {
	Parameters []string // The parameter types you want to look for (e.g. string, string)
}

func (f *FilterGoFileFunctionParameter) Scan(content string) (bool, string) {
	trimmed := strings.TrimSpace(content)

	// Search for functions
	if !strings.HasPrefix(trimmed, "func") {
		return false, ""
	}

	// Extract the function arguments by getting everything inside the brackets
	arguments := strings.TrimFunc(trimmed, func(r rune) bool {
		return r != '(' && r != ')'
	})

	// Get the name of the function for later
	fnName := strings.TrimRightFunc(trimmed, func(r rune) bool {
		return r != '('
	})
	fnName = strings.TrimLeftFunc(fnName, func(r rune) bool {
		return !unicode.IsSpace(r)
	})
	fnName = strings.TrimSpace(fnName)
	fnName = strings.Trim(fnName, "()")

	// Go through all parameters and extract the types
	parameters := []string{}
	for param := range strings.SplitSeq(arguments, ",") {
		param = strings.TrimSpace(param)
		param = strings.Trim(param, "()")
		param = strings.TrimLeftFunc(param, func(r rune) bool {
			return !unicode.IsSpace(r)
		})
		if param == "" {
			continue
		}
		parameters = append(parameters, strings.TrimSpace(param))
	}

	// If they match, return the result
	if slices.Equal(parameters, f.Parameters) {
		return true, fnName
	}
	return false, ""
}

// Get the package name in the file
var FilterGoFilePackageName = filterGoFilePackageName{}

type filterGoFilePackageName struct{}

func (f filterGoFilePackageName) Scan(content string) (bool, string) {

	// Make sure the line starts with package
	if !strings.HasPrefix(content, "package") {
		return false, ""
	}
	content, _ = strings.CutPrefix(content, "package")
	content = strings.TrimSpace(content)

	// Make sure what's remaining is truely the package
	if strings.ContainsFunc(content, unicode.IsSpace) || content == "" {
		return false, ""
	}

	return true, content
}
