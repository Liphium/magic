package mrunner

import "strings"

// Helper function for removing // comments in a line
func RemoveOneLineComments(line string) string {
	args := strings.SplitN(line, "//", 2)
	return args[0]
}
