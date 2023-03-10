#!/bin/bash
set -eu -o pipefail

find_banned_imports () {
  go list -f '{{ join .Deps "\n" }}' "$@" |
  grep -E -f "$(dirname "$0")/test-banned-imports.txt"
}

if [ $# -ne 0 ]; then
  echo "usage: $0" >&2
  exit 2
fi

if found=$(find_banned_imports github.com/rkoesters/xkcd-gtk/cmd/xkcd-gtk); then
    printf 'ERROR: found banned import(s):\n%s\n' "${found:?}" >&2
    exit 1
fi

echo "$0" "PASSED" >&2
