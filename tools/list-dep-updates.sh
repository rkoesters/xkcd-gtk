#!/bin/sh
# Find and print updates to dependencies go modules.
set -eu

dep_modules () {
	go mod graph |
	grep '^github.com/rkoesters/xkcd-gtk' |
	cut -d ' ' -f 2 |
	sort |
	uniq
}

go list -m -u $(dep_modules) 2>/dev/null |
grep -e '\[.*\]'
