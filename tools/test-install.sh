#!/bin/sh
# Comfirms that `make uninstall` removes everything that `make install` creates.
set -eu

echo "mkdir out"
mkdir out
make prefix=out install
make prefix=out uninstall
echo "rmdir out/..."
rmdir out/*/*/*/*/* out/*/*/*/* out/*/*/* out/*/* out/* out
