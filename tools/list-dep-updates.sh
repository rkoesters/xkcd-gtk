#!/bin/sh
# Find and print updates to dependencies go modules.
set -eu

go list -m -u $(tools/list-deps.sh) 2>/dev/null |
grep -e '\[.*\]'
