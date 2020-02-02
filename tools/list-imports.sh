#!/bin/sh
go list -f '{{ join .Imports "\n" }}' "$@" | # List all of our imports.
grep '\.' | # Filter out all imports that don't appear to be URLs.
grep -v 'github.com/rkoesters/xkcd-gtk' | # Filter out ourselves.
sort | # Sort imports so 'uniq' can work properly.
uniq # Filter out repeated imports.
