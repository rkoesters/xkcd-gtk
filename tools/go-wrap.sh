#!/bin/sh
# Simple script to wrap a plain text file as go const string.
set -eu

if [ $# != 1 ]; then
	echo "usage: $0 FILE"
	exit 1
fi

printf 'package main\n\nconst '

basename "$1" | tr '.-' '_' | tr -d '\n' |
sed -e 's/_\([^_]*\)$/\U\1/g' -e 's/_\([a-z]\)/\U\1/g'

printf ' = `'

cat "$1"

printf '`\n'
