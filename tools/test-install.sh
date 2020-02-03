#!/bin/bash
# Comfirms that `make uninstall` removes everything that `make install` creates.
set -eEu
trap 'echo "$0 FAILED"' ERR

verbose=false

while [[ $# > 0 ]]; do
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

print_and_run mkdir out
print_and_make install prefix=out
print_and_make uninstall prefix=out
print_and_run rmdir out/*/*/*/*/* out/*/*/*/* out/*/*/* out/*/* out/* out

echo "$0 PASSED"
