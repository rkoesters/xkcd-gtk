#!/bin/sh
# Find and print versions (or commits) of the dependencies of the given Go
# packages.
set -eu

get_package_version () {
	git -C "$1" describe --always --tags --dirty
}

print_package_version () {
	printf '%s\t%s\n' "$1" "$(get_package_version "$1")"
}

# Run all our tools/* before we change directory.
deps="$(tools/list-deps.sh "$@")"

# Change directory to $GOPATH/src so that package names work as relative paths
# to the respective package.
cd "$(go env GOPATH)/src"

for package in $deps; do
	print_package_version "$package"
done
