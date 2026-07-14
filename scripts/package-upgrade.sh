#!/bin/bash
find "$(git rev-parse --show-toplevel)" -name "go.mod" -execdir sh -c "go get -u ./... && go mod tidy" \;
