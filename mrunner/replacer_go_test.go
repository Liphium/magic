package mrunner_test

import (
	"strings"
	"testing"

	"github.com/Liphium/magic/mrunner"
)

func TestCommentCleaner(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no comments",
			input:    "line1\nline2",
			expected: "line1\nline2",
		},
		{
			name:     "single line comment",
			input:    "line1 // comment\nline2",
			expected: "line1 \nline2",
		},
		{
			name:     "multiple single line comments",
			input:    "line1 // comment1\nline2 // comment2",
			expected: "line1 \nline2 ",
		},
		{
			name:     "block comment",
			input:    "line1 /* comment */\nline2",
			expected: "line1 \nline2",
		},
		{
			name:     "block comment spanning multiple lines",
			input:    "line1 /* comment\nstill comment */\nline2",
			expected: "line1 \nline2",
		},
		{
			name:     "mixed comments",
			input:    "line1 // single line\nline2 /* block \n comment */ line3 // another single",
			expected: "line1 \nline2 \n line3 ",
		},
		{
			name:     "comment at beginning of line",
			input:    "// comment\nline1",
			expected: "line1",
		},
		{
			name:     "block comment at beginning of line",
			input:    "/* comment */line1",
			expected: "line1",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "input with only comments",
			input:    "// line1\n/* block \n comment */",
			expected: "",
		},
		{
			name:     "complex nested and mixed",
			input:    "1 /* 2 3 /* \n 4 \n */ 5 /* 6 */",
			expected: "1 \n 5 ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleaner := mrunner.CommentCleaner{}
			actual := mrunner.ReplaceLines(tt.input, &cleaner)
			if actual != tt.expected && strings.Trim(actual, "\n") != tt.expected {
				t.Errorf("with input %q expected %q, got %q", tt.input, tt.expected, actual)
			}
		})
	}
}

func TestPackageReplacer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		newPkg   string
		expected string
	}{
		{
			name:     "replace package name",
			input:    "package oldname",
			newPkg:   "newname",
			expected: "package newname",
		},
		{
			name:     "no package line",
			input:    "import \"fmt\"",
			newPkg:   "main",
			expected: "import \"fmt\"",
		},
		{
			name:     "package line with spaces",
			input:    "   package   old",
			newPkg:   "foo",
			expected: "package foo",
		},
		{
			name:     "empty input",
			input:    "",
			newPkg:   "bar",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replacer := mrunner.GoPackageReplacer{NewPackage: tt.newPkg}
			actual := replacer.Replace(tt.input)
			if actual != tt.expected {
				t.Errorf("with input %q and newPkg %q expected %q, got %q", tt.input, tt.newPkg, tt.expected, actual)
			}
		})
	}
}
