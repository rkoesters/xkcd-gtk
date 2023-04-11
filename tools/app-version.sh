#!/bin/bash
set -eu -o pipefail

version_from_git () {
  git describe --tags --match='v[0-9].[0-9]*.[0-9]*' --dirty
}

version_from_appdata () {
  grep '<release version="' data/com.github.rkoesters.xkcd-gtk.appdata.xml |
  head -n 1 |
  cut -d '"' -f 2
}

version_from_ci () {
  if [ "${CI:-false}" == "false" ]; then
    return 1
  fi
  printf 'ci-%s-%s\n' "${GITHUB_REF_NAME:-nullref}" "${GITHUB_SHA:-nullsha}"
}

version_from_git 2>/dev/null ||
version_from_ci 2>/dev/null ||
version_from_appdata 2>/dev/null ||
echo "unknown"
