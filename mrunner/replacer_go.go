package mrunner

import (
	"strings"
)

// Removes all one line comments in a file (with //)
type OneLineCommentReplacer struct{}

func (r OneLineCommentReplacer) Replace(old string) string {
	return RemoveOneLineComments(old)
}

// Removes all block comments and one line comments (/* */ and //)
type CommentCleaner struct {
	started bool
}

func (r *CommentCleaner) Replace(old string) string {
	old = RemoveOneLineComments(old)

	// Remove everything that's in block comments
	new := ""
	startCount := 0
	for start := range strings.SplitSeq(old, "/*") {
		startCount++

		// If we're inside of a block comment, look for the end and mark as started
		if startCount >= 2 || r.started {
			r.started = true

			// Look for the end of block comments
			endCount := 0
			for end := range strings.SplitSeq(start, "*/") {
				endCount++

				// If we found an end, mark as ended and add to new string
				if endCount%2 == 0 || !r.started {
					r.started = false
					new += end
					continue
				}

				r.started = true
			}
			continue
		}

		// Add in case we're behind a block comment start sequence
		r.started = false
		new += start
	}

	return new
}

type GoPackageReplacer struct {
	NewPackage string // New package name instead of the old one
}

func (r GoPackageReplacer) Replace(old string) string {
	trimmed := strings.TrimSpace(old)

	// Check for the package line
	if strings.HasPrefix(trimmed, "package") {
		return "package " + r.NewPackage
	}

	return old
}
