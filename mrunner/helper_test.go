package mrunner_test

import (
	"testing"

	"github.com/Liphium/magic/mrunner"
)

func TestSingleLineComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "line with comment",
			input:    "some code // this is a comment",
			expected: "some code ",
		},
		{
			name:     "line without comment",
			input:    "some code without a comment",
			expected: "some code without a comment",
		},
		{
			name:     "empty line",
			input:    "",
			expected: "",
		},
		{
			name:     "line with only comment",
			input:    "// only a comment",
			expected: "",
		},
		{
			name:     "line with multiple slashes but not a comment",
			input:    "http://example.com",
			expected: "http:",
		},
		{
			name:     "line with comment and leading/trailing spaces",
			input:    "  some code // comment  ",
			expected: "  some code ",
		},
		{
			name:     "line with no space before comment",
			input:    "code//comment",
			expected: "code",
		},
		{
			name:     "line with three slashes",
			input:    "code /// comment",
			expected: "code ",
		},
		{
			name:     "line ending with single slash",
			input:    "code /",
			expected: "code /",
		},
		{
			name:     "double comments",
			input:    "code // hello world // more comments",
			expected: "code ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := mrunner.RemoveOneLineComments(tt.input)
			if actual != tt.expected {
				t.Errorf("with input %q expected %q, got %q", tt.input, tt.expected, actual)
			}
		})
	}
}
