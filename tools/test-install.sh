#!/bin/sh
# Comfirms that `make uninstall` removes everything that `make install` creates.
set -eu

echo "mkdir out"
mkdir out
make install prefix=out
make uninstall prefix=out
echo "rmdir out/..."
rmdir out/*/*/*/*/* out/*/*/*/* out/*/*/* out/*/* out/* out
