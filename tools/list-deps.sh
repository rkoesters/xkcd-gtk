#!/bin/sh
# Find and print the dependencies of this go module.
set -eu

go mod graph |
sed 's/^.* //g' |
sort |
uniq
