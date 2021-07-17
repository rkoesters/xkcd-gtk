#!/bin/sh
# Find and print updates to direct dependencies of the current module.
set -eu

go list -m -u $(tools/list-mod-deps.sh) 2>/dev/null |
grep -e '\[.*\]'
