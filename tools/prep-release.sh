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

echo "Validating ${appdata_xml:?}"
if ! appstream-util validate-relax "${appdata_xml:?}"; then
  failure "invalid appdata"
fi

echo "Checking for correct attribution in ${appdata_xml:?}"
if grep -v '@2x.png' "${appdata_xml:?}" | grep -q '@'; then
  failure "appdata includes '@' symbol which doesn't work outside GitHub"
fi

echo "Reading most recent release from changelog"
if ! release=$(grep "<release version=" "${appdata_xml:?}" | head -n 1 | cut -d '"' -f 2); then
  failure "could not read most recent release from changelog"
fi

echo "Checking if current commit is tagged ${release:?}"
if [ "${release:?}" != "$(git describe --exact-match --tags --match='[0-9].[0-9]*.[0-9]*')" ]; then
  failure "current commit not tagged ${release:?}"
fi

echo "Checking for date for release ${release:?} in changelog"
date=$(git log -1 --format='%ad' --date=short -- "${release:?}")
if ! grep -q "<release version=\"${release:?}\" date=\"${date:?}\"" "${appdata_xml:?}"; then
  failure "date ${date:?} not found in appdata changelog for version ${release:?}"
fi

success "checks passed!"
