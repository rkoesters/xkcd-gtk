#!/bin/sh
# Find and print modules that are direct dependencies of the current module.
set -eu

go mod graph |
grep -E '^github.com/rkoesters/xkcd-gtk ' |
cut -d ' ' -f 2 |
grep -F "$(tools/list-pkg-deps.sh github.com/rkoesters/xkcd-gtk/cmd/xkcd-gtk)" |
sort |
uniq
