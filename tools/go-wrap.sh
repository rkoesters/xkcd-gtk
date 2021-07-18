#!/bin/sh
# Simple script to wrap a plain text file as go const string.
set -eu

if [ $# != 1 ]; then
  echo "usage: $0 FILE" >&2
  exit 1
fi

file="$1"
dir="$(dirname "$file")"
case "$dir" in
  *cmd/*) package="main" ;;
  *) package="$(basename "$dir")" ;;
esac

printf 'package %s\n\nconst ' "$package"

basename "$file" | tr '.-' '_' | tr -d '\n' |
sed -e 's/_\([^_]*\)$/\U\1/g' -e 's/_\([a-z]\)/\U\1/g'

printf ' = `'

cat "$file"

printf '`\n'
