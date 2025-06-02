package mrunner

import (
	"strings"
)

// Helper function for removing // comments in a line
func RemoveOneLineComments(line string) string {
	newLine := ""
	quoteCount := 0
	for text := range strings.SplitSeq(line, "\"") {
		// If it ends in backslash, end the quote
		if strings.HasSuffix(text, "\\") {
			newLine += text + "\""
			continue
		}
		quoteCount++

		if quoteCount%2 == 0 {
			newLine += text + "\""
			continue
		}

		// Check if a comment started
		args := strings.SplitN(text, "//", 2)
		if len(args) > 1 {
			if newLine != "" {
				newLine, _ = strings.CutSuffix(newLine, "\"")
				return newLine + "\"" + args[0]
			}
			return args[0]
		}
		newLine += text + "\""
	}

	newLine, _ = strings.CutSuffix(newLine, "\"")
	return newLine
}
