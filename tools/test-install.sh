#!/bin/sh
# Comfirms that `make uninstall` removes everything that `make install` creates.
set -eu

echo "mkdir out"
mkdir out
make --no-print-directory install prefix=out
make --no-print-directory uninstall prefix=out
echo "rmdir out/..."
rmdir out/*/*/*/*/* out/*/*/*/* out/*/*/* out/*/* out/* out
