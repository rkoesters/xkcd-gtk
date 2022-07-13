#!/bin/bash
set -eu -o pipefail

pkg-config --modversion pango |
sed 's/\.[^.]*$//g' |
sed 's/\./_/g' |
sed 's/^/pango_/g'
