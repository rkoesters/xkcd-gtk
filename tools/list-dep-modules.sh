#!/bin/sh
# Find and print modules on which this module depends.
set -eu

go mod graph |
grep '^github.com/rkoesters/xkcd-gtk' |
cut -d ' ' -f 2 |
sort |
uniq
