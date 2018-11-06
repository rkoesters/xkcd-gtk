#!/bin/sh -eu
appdata_xml="data/com.github.rkoesters.xkcd-gtk.appdata.xml"

echo "Validating appdata..."
appstream-util validate-relax "$appdata_xml"

echo "Finding git tag..."
version=$(git describe --exact-match --tags HEAD)

echo "Checking for git tag in appdata changelog..."
if ! grep -q "<release version=\"$version\"" "$appdata_xml"; then
	echo "version $version not found in appdata changelog"
	exit 1
fi

echo "Checks passed!"
