#!/bin/bash
# Find and print modules that are direct and indirect dependencies of the this
# application.
set -eu -o pipefail

if ! pkg_deps="$(tools/list-pkg-deps.sh github.com/rkoesters/xkcd-gtk/cmd/xkcd-gtk)"; then
  echo "$0: error finding dependencies of cmd/xkcd-gtk" >&2
  exit 1
fi

go mod graph |
cut -d ' ' -f 2 |
grep -F "${pkg_deps:?}" |
sort |
uniq
