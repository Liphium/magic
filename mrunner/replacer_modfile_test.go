package mrunner_test

import (
	"testing"

	"github.com/Liphium/magic/mrunner"
)

func TestVersionReplacer(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		version string
		want    string
	}{
		{
			name:    "replace go version",
			input:   "go 1.18",
			version: "1.21",
			want:    "go 1.21",
		},
		{
			name:    "replace go version with extra spaces",
			input:   "  go 1.19  ",
			version: "1.20",
			want:    "go 1.20",
		},
		{
			name:    "do not replace unrelated line",
			input:   "module github.com/example/project",
			version: "1.22",
			want:    "module github.com/example/project",
		},
		{
			name:    "do not replace empty line",
			input:   " ",
			version: "1.23",
			want:    " ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mrunner.GoVersionReplacer{Version: tt.version}
			got := r.Replace(tt.input)
			if got != tt.want {
				t.Errorf("Replace(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
