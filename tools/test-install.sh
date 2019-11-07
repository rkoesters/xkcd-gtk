#!/bin/sh
set -eu

mkdir out
make prefix=out install
make prefix=out uninstall
rmdir out/*/*/*/*/* out/*/*/*/* out/*/*/* out/*/* out/* out
