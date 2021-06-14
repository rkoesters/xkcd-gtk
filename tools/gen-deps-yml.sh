#!/bin/sh
set -eu

IFS='
'
for dep in $(tools/list-dep-modules.sh); do
	name=$(echo "${dep:?}" | cut -d '@' -f 1)
	version=$(echo "${dep:?}" | cut -d '@' -f 2)

	echo "      - type: git"
	echo "        url: https://${name:?}.git"
	echo "        tag: ${version:?}"
	echo "        dest: vendor/${name:?}"
	echo
done
