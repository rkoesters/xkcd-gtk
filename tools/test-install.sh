#!/bin/bash
# Comfirms that `make uninstall` removes everything that `make install` creates.
set -eEu
trap 'echo "$0 FAILED"' ERR

export MAKEFLAGS=""

verbose=false

while [ $# -gt 0 ]; do
  case "$1" in
    -v|--verbose)
      verbose=true
      ;;
    *)
      echo "usage: $0 [ -v | --verbose ]" >&2
      exit 1
      ;;
  esac
  shift
done

print_and_run () {
  if [[ $verbose == true ]]; then
    echo "$*"
  fi
  "$@"
}

print_and_make () {
  if [[ $verbose = true ]]; then
    print_and_run make "$@"
  else
    print_and_run make -s "$@"
  fi
}

tmpdir=$(mktemp -d /dev/shm/xkcd-gtk.XXXXXXXX)

print_and_make install DESTDIR="$tmpdir"
print_and_make uninstall DESTDIR="$tmpdir"
print_and_run rmdir "$tmpdir"/*/*/*/*/*/*
print_and_run rmdir "$tmpdir"/*/*/*/*/*
print_and_run rmdir "$tmpdir"/*/*/*/*
print_and_run rmdir "$tmpdir"/*/*/*
print_and_run rmdir "$tmpdir"/*/*
print_and_run rmdir "$tmpdir"/*
print_and_run rmdir "$tmpdir"

echo "$0 PASSED"
