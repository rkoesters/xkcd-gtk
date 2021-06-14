#!/bin/sh
# Find and print modules on which this module depends.
set -eu

go mod graph |
cut -d ' ' -f 2 |
grep -F "$(tools/list-deps.sh github.com/rkoesters/xkcd-gtk/cmd/xkcd-gtk)" |
sort |
uniq
