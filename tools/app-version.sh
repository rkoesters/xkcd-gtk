#!/bin/bash
set -eu -o pipefail

git describe --tags --match='[0-9].[0-9]*.[0-9]*' --dirty 2>/dev/null ||
grep '<release version="' data/com.github.rkoesters.xkcd-gtk.appdata.xml |
  head -n 1 | cut -d '"' -f 2 ||
echo "unknown"
