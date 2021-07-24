#!/bin/sh
# Confirm that there are no conflicting source directories in the flatpak
# config.
if [ $# != 1 ]; then
  echo "usage: $0 flatpak-yml-file" >&2
  exit 2
fi
yml="${1:?}"

if [ "$(grep dest: "${yml:?}" | uniq -D)" != "" ]; then
  echo "ERROR: duplicate 'dest:' in flatpak config:" >&2
  duplicates=$(
    grep dest: "${yml:?}" |
    sed -e "s/^ *dest: //g" -e "s/ *$//g" |
    sort |
    uniq -d
  )
  for dup in ${duplicates:?}; do
    grep -n dest: "${yml:?}" | grep "${dup:?}"
  done
  echo "$0 FAILED"
  exit 1
fi

echo "$0 $yml PASSED"
