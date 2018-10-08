#!/bin/sh
git describe --always --tags --dirty 2>/dev/null ||
dpkg-parsechangelog -S Version 2>/dev/null ||
echo "unknown"
