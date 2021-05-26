#!/bin/sh
pkg-config --modversion pango |
sed 's/\.[^.]*$//g' |
sed 's/\./_/g' |
sed 's/^/pango_/g'
