#!/bin/sh -eu
appstream-util validate-relax data/*.appdata.xml
git describe --exact-match --tags HEAD
