#!/bin/bash
find "$(git rev-parse --show-toplevel)" -name "go.mod" -execdir go mod tidy \;
