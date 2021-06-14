#!/bin/sh
# Find and print the dependencies of this go module.
set -eu

go mod graph |
cut -d ' ' -f 2 |
sort |
uniq
