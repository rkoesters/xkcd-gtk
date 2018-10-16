#!/bin/sh -eu
kernel="$(uname -s | tr 'A-Z' 'a-z')"
processor="$(uname -m | sed -e 's/x86_64/amd64/g')"

package_url="https://dl.google.com/go/go1.10.4.$kernel-$processor.tar.gz"

curl -s "$package_url" | tar -xz

GO="$(pwd)/go/bin/go"

if [ -x "$GO" ]; then
	echo "$GO"
else
	echo "failed to bootstrap: cannot execute '$GO'" >&2
	exit 1
fi
