#!/bin/sh -eu
appdata_xml="data/com.github.rkoesters.xkcd-gtk.appdata.xml"

appstream-util validate-relax "$appdata_xml"

version=$(git describe --exact-match --tags HEAD)

if ! grep -q "<release version=\"$version\"" "$appdata_xml"; then
	echo "version $version not found in appdata changelog"
	exit 1
fi
