#!/bin/sh
# Confirm that there are no conflicting source directories in the flatpak
# config.
yml=com.github.rkoesters.xkcd-gtk.yml

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

echo "$0 PASSED"
