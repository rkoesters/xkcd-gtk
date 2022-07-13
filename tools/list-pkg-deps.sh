#!/bin/bash
# Find and print the dependencies of the given go packages.
set -eu -o pipefail

if [ "$(go version | cut -d ' ' -f 3 | cut -d '.' -f 2)" -ge "18" ]; then
  go list -buildvcs=false -f '{{ join .Deps "\n" }}' "$@"
else
  go list -f '{{ join .Deps "\n" }}' "$@"
fi |
grep -v -e '^internal/' -e '^vendor/' -e '^github.com/rkoesters/xkcd-gtk' |
grep '^.*\..*/' |
sed -e 's#\([^/]*/[^/]*/[^/]*\).*#\1#g' |
sort |
uniq
