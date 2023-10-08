#!/bin/bash
set -eu -o pipefail

pkg_config_tag () {
    pkg-config --modversion "${2:?}" |
    sed 's/\.[^.]*$//g' |
    sed 's/\./_/g' |
    sed "s/^/${1:?}_/g"
}

pkg_config_tag cairo 'cairo'
pkg_config_tag gdk_pixbuf 'gdk-pixbuf-2.0'
pkg_config_tag glib 'glib-2.0'
pkg_config_tag gtk 'gtk+-3.0'
pkg_config_tag pango 'pango'
