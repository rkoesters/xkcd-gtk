#!/bin/sh
# Find and print updates to direct and indirect dependencies of the current
# module.
set -eu

# NOTE(SC2046): The word splitting is intentional.
# shellcheck disable=SC2046
go list -m -u $(tools/list-all-mod-deps.sh) 2>/dev/null |
grep -e '\[.*\]'
