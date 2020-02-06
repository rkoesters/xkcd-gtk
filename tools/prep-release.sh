#!/bin/sh
set -eu

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

echo "Checking for correct attribution in $appdata_xml"
if grep -q '@' "$appdata_xml"; then
	failure "appdata includes '@' symbol which doesn't work outside GitHub"
fi

echo "Checking for tag matching current commit"
tag=$(git describe --exact-match --tags)
if [ $? != 0 ]; then
	failure "could not find tag for current commit"
fi

echo "Checking for tag $tag in changelog"
if ! grep -q "<release version=\"$tag\"" "$appdata_xml"; then
	failure "version $tag not found in appdata changelog"
fi

echo "Checking for date for tag $tag in changelog"
date=$(git log -1 --format='%ad' --date=short "$tag")
if ! grep -q "<release version=\"$tag\" date=\"$date\"" "$appdata_xml"; then
	failure "date $date not found in appdata changelog for version $tag"
fi

success "checks passed!"
