#!/bin/sh -eu
# simple script to wrap a plain text file as go const string

if [ $# != 1 ]; then
	echo "usage: $0 FILE"
	exit 1
fi

constName="$(basename "$1" | tr '.-' '_')"
constValue="$(cat "$1")"

printf 'package main\n\nconst %s = `%s\n`\n' "$constName" "$constValue"
