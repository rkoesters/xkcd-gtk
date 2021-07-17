#!/bin/sh
# Confirm that there are no conflicting source directories in the flatpak
# config.
if [ "$(grep dest: com.github.rkoesters.xkcd-gtk.yml | uniq -D)" != "" ]; then
	echo "ERROR: duplicate 'dest:' in flatpak config:" >&2
	duplicates=$(
		grep dest: com.github.rkoesters.xkcd-gtk.yml |
		sed -e "s/^ *dest: //g" -e "s/ *$//g" |
		sort |
		uniq -d
	)
	for dup in ${duplicates:?}; do
		grep -n dest: com.github.rkoesters.xkcd-gtk.yml | grep "${dup:?}"
	done
	echo "$0 FAILED"
	exit 1
fi

echo "$0 PASSED"
