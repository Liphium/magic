package integration

import "regexp"

func IsPathSanitized(path string) bool {
	pathRegex := "^[A-Za-z_][A-Za-z0-9_]*$"
	found, err := regexp.MatchString(pathRegex, path)
	if err != nil {
		return false
	}
	return found
}
