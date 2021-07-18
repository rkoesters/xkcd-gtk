#!/bin/sh
# Find and print modules that are direct and indirect dependencies of the this
# application.
set -eu

go mod graph |
cut -d ' ' -f 2 |
grep -F "$(tools/list-pkg-deps.sh github.com/rkoesters/xkcd-gtk/cmd/xkcd-gtk)" |
sort |
uniq
