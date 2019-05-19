#!/bin/sh -u
failure () {
	echo "FAILURE $*"
	exit 1
}
success () {
	echo "SUCCESS $*"
}

appdata_xml="data/com.github.rkoesters.xkcd-gtk.appdata.xml.in"

echo "Validating $appdata_xml"
if ! appstream-util validate-relax "$appdata_xml"; then
	failure "invalid appdata"
fi

echo "Checking for tag matching current commit"
version=$(git describe --exact-match --tags)
if [ $? != 0 ]; then
	failure "could not find tag for current commit"
fi

echo "Checking for tag in changelog"
if ! grep -q "<release version=\"$version\"" "$appdata_xml"; then
	failure "version $version not found in appdata changelog"
fi

echo "Checking for date for tag in changelog"
date=$(git log -1 --format='%ad' --date=short "$version")
if ! grep -q "<release version=\"$version\" date=\"$date\"" "$appdata_xml"; then
	failure "date $date not found in appdata changelog for version $version"
fi

success "checks passed!"
