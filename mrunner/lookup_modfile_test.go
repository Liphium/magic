package mrunner_test

import (
	"testing"

	"github.com/Liphium/magic/mrunner"
)

func TestReplacersFilter(t *testing.T) {
	tests := []struct {
		line        string
		expected    string
		shouldMatch bool
	}{
		{"replace github.com/old/module => github.com/new/module v1.2.3", "github.com/old/module;github.com/new/module v1.2.3", true},
		{"replace github.com/old/module v1.0.0 => github.com/new/module v1.2.3", "github.com/old/module v1.0.0;github.com/new/module v1.2.3", true},
		{"replace example.com/foo => ../local/foo", "example.com/foo;../local/foo", true},
		{"replace something => else", "something;else", true},
		{"notreplace github.com/old/module => github.com/new/module", "", false},
		{"replace github.com/old/module github.com/new/module", "", false},
	}

	for _, tt := range tests {
		matched, result := mrunner.FilterModFileReplacers.Scan(tt.line)
		if matched != tt.shouldMatch {
			t.Errorf("For line '%s', expected match=%v but got %v", tt.line, tt.shouldMatch, matched)
		}
		if matched && result != tt.expected {
			t.Errorf("For line '%s', expected result '%s' but got '%s'", tt.line, tt.expected, result)
		}
	}
}
