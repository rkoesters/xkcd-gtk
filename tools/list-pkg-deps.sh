#!/bin/sh
# Find and print the dependencies of the given go packages.
set -eu

go list -f '{{ join .Deps "\n" }}' "$@" |
grep -v -e '^internal/' -e '^vendor/' -e '^github.com/rkoesters/xkcd-gtk' |
grep '^.*\..*/' |
sed -e 's#\([^/]*/[^/]*/[^/]*\).*#\1#g' |
sort |
uniq
