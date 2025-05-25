package integration

import "strings"

func SnakeToCamelCase(s string, capitalizeFirst bool) string {
	// Determine the start for the capitialization
	start := 1
	if capitalizeFirst {
		start = 0
	}

	// Start converting and return result
	parts := strings.Split(s, "_")
	for i := start; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][0:1]) + strings.ToLower(parts[i][1:])
	}
	return strings.Join(parts, "")
}
