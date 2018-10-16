#!/bin/sh
go list -f '{{ join .Deps "\n" }}' "$@" | # List all of our dependencies.
grep '\.' | # Filter out all dependencies that don't appear to be URLs.
sort | # Sort dependencies so 'uniq' can work properly.
uniq # Filter out repeated dependencies.
