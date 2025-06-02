package mrunner

import "strings"

// Replace the go version in a go.mod or go.work file
type GoVersionReplacer struct {
	Version string // Version to replace the go version with
}

func (r *GoVersionReplacer) Replace(old string) string {
	trimmed := strings.TrimSpace(old)
	if strings.HasPrefix(trimmed, "go") {
		return "go " + r.Version
	}
	return old
}
