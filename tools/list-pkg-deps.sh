#!/bin/bash
# Find and print the dependencies of the given go packages.
set -eu -o pipefail

go list -buildvcs=false -f '{{ join .Deps "\n" }}' "$@" |
grep -v '^github.com/rkoesters/xkcd-gtk' |
grep '^[^/]*\.[^/]*/' |
sed -e 's#\([^/]*/[^/]*/[^/]*\).*#\1#g' |
sort |
uniq
