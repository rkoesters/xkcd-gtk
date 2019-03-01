#!/bin/sh -eu
# simple script to wrap a plain text file as go const string

if [ $# != 1 ]; then
	echo "usage: $0 FILE"
	exit 1
fi

constName="$(basename "$1" | tr '.-' '_' | sed -e 's/_\([^_]*\)$/\U\1/g' -e 's/_\([a-z]\)/\U\1/g')"
constValue="$(cat "$1")"

printf 'package main\n\nconst %s = `%s\n`\n' "$constName" "$constValue"
