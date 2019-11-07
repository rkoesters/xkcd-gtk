#!/bin/sh
# Find and print versions (or commits) of the dependencies of the given Go
# packages.
set -eu

version () {
	git -C "$1" describe --always --tags --dirty
}

# Run all our tools/* before we change directory.
deps="$(tools/list-deps.sh "$@")"

# Change directory to $GOPATH/src so that package names work as relative paths
# to the respective package.
cd "$(go env GOPATH)/src"

(for pkg in $deps; do echo "$pkg" "$(version "$pkg")"; done) | column -t -s ' '
