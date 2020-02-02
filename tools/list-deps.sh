#!/bin/sh
go list -f '{{ join .Deps "\n" }}' "$@" | # List all of our dependencies.
grep '\.' | # Filter out all dependencies that don't appear to be URLs.
sed 's#^\([^/]\+/[^/]\+/[^/]\+\).*$#\1#g' | # Trim URLs to repo root (e.g. github.com/user/repo).
grep -v 'github.com/rkoesters/xkcd-gtk' | # Filter out ourselves.
sort | # Sort dependencies so 'uniq' can work properly.
uniq # Filter out repeated dependencies.
