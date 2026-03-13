package mservices_test

import (
	"testing"

	mservices "github.com/Liphium/magic/v3/mrunner/services"
)

func TestMajorVersionGetting(t *testing.T) {
	tests := []struct {
		image    string
		expected int
	}{
		{"postgres:17", 17},
		{"node:v20.1.0", 20},
		{"nginx:1.25-alpine", 1},
		{"redis:7.2.4", 7},
		{"ubuntu:22.04", 22},
		{"alpine:3.18.4", 3},
		{"myimage:latest", -1},
		{"myimage", -1},
		{"myimage:v2", 2},
		{"myimage:2", 2},
		{"registry.example.com/myimage:5.3", 5},
		{"registry.example.com:5000/myimage:10.1", 10},
		{"noversion:", -1},
		{"", -1},
	}

	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			got := mservices.GetImageMajorVersion(tt.image)
			if got != tt.expected {
				t.Errorf("getImageMajorVersion(%q) = %d, want %d", tt.image, got, tt.expected)
			}
		})
	}
}
