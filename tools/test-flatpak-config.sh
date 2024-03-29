#!/bin/bash
set -eu -o pipefail

# Confirm that there are no conflicting source directories in the flatpak
# config.
if [ $# -lt 1 ]; then
  echo "usage: $0 flatpak-yml-file [ flatpak-yml-file ... ]" >&2
  exit 2
fi

while [ $# -gt 0 ]; do
  yml="${1:?}"
  shift

  if [ ! -r "${yml:?}" ]; then
    echo 'ERROR: cannot read flatpak config:' "${yml:?}" >&2
    exit 1
  fi

  if [ ! -f "${yml:?}" ]; then
    echo 'ERROR: flatpak config is not a regular file:' "${yml:?}" >&2
    exit 1
  fi

  if [ "$(grep dest: "${yml:?}" | grep -v -e 'dest: go$' |
          sort | uniq -D)" != "" ]; then
    echo "ERROR: duplicate 'dest:' in flatpak config:" >&2
    duplicates=$(
      grep dest: "${yml:?}" |
      grep -v -e 'dest: go$' |
      sed -e "s/^ *dest: //g" -e "s/ *$//g" |
      sort |
      uniq -d
    )
    for dup in ${duplicates:?}; do
      grep -n dest: "${yml:?}" | grep "${dup:?}"
    done
    exit 1
  fi

  echo "$0" "${yml:?}" "PASSED" >&2
done
