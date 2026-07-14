#!/bin/bash
NEW_VER=$1
if [ -z "$NEW_VER" ]; then
    echo "Need version. Usage: $0 1.22"
    exit 1
fi
find "$(git rev-parse --show-toplevel)" -name "go.mod" -execdir go mod edit -go="$NEW_VER" \;
