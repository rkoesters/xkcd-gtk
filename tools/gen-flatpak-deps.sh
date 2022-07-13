#!/bin/sh
set -eu

LC_COLLATE=C
export LC_COLLATE

name_to_remote () {
  printf 'https://'
  printf '%s' "${1:?}" |
  sed \
    -e 's#/v[0-9][0-9]*$##g' \
    -e 's#go.etcd.io/#github.com/etcd-io/#g' \
    -e 's#golang.org/x/#go.googlesource.com/#g' \
  # This comment ensures printf does not become an argument to sed.
  printf '.git\n'
}

list_deps () {
  tools/list-app-deps.sh |
  tr '@' ' ' |
  sort -r |
  rev |
  uniq -f 1 |
  rev |
  sort
}

deps="$(list_deps)"
if ! [ $? ]; then
  echo "$0: error generating list of dependencies" >&2
  exit 1
fi

IFS='
'
for dep in ${deps:?}; do
  name=$(echo "${dep:?}" | cut -d ' ' -f 1)
  version=$(echo "${dep:?}" | cut -d ' ' -f 2)
  remote="$(name_to_remote "${name:?}")"

  echo
  echo "      - type: git"
  echo "        url: ${remote:?}"
  case ${version:?} in
    v0.0.0-*)
      commit="$(echo "${version:?}" | cut -d '-' -f 3)"
      echo "        # ${version:?}"
      echo "        commit: ${commit:?}"
      ;;
    *)
      echo "        tag: ${version:?}"
      ;;
  esac
  echo "        dest: vendor/${name:?}"
done
