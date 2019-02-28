#!/bin/sh -eu
# simple script to wrap plain text file as go const string

if [ $# != 1 ]; then
	echo "usage: $0 FILE"
	exit 1
fi

varName="$(basename "$1" '.in' | tr '.-' '_')"

printf 'package main\n\nconst %s = `%s`\n' "$varName" "$(cat "$1")"
