#!/bin/sh
# Find and print versions of the dependencies of all go.mod files under the
# current directory. Run `make vendor` before running to get versions of all
# dependencies.
set -eu

cat $(find . -name 'go.mod' -type f) |
grep '	' |
tr -d '\t' |
sort |
uniq
