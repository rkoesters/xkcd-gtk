#!/bin/sh
# Comfirms that `make uninstall` removes everything that `make install` creates.
set -eu

print_and_run () {
	echo "$*"
	"$@"
}

print_and_run mkdir out
print_and_run make install prefix=out
print_and_run make uninstall prefix=out
print_and_run rmdir out/*/*/*/*/* out/*/*/*/* out/*/*/* out/*/* out/* out
