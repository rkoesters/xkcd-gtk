#!/bin/sh
go list -f '{{ join .Deps "\n" }}' "$@" | # List all of our dependencies.
grep '^[^/]\+\.[^/]\+/[^/]\+/[^/]\+$' | # Filter to repos (e.g. github.com/*/*).
sort | # Sort dependencies so 'uniq' can work properly.
uniq # Filter out repeated dependencies.
