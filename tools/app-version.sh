#!/bin/sh
git describe --always --tags --dirty 2>/dev/null ||
dpkg-parsechangelog -S Version 2>/dev/null ||
basename $PWD | grep -o '[0-9][0-9]*\.[0-9][0-9]*\.[0-9].*' ||
echo "unknown"
