#!/bin/sh
set -eu

echo "mkdir out"
mkdir out
make prefix=out install
make prefix=out uninstall
echo "rmdir out/..."
rmdir out/*/*/*/*/* out/*/*/*/* out/*/*/* out/*/* out/* out
