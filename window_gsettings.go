// +build linux freebsd netbsd openbsd

package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func lookupHeaderBarIconSize() gtk.IconSize {
	interfaceSettings := glib.SettingsNew("org.gnome.desktop.interface")
	theme := interfaceSettings.GetString("gtk-theme")
	if theme == "elementary" {
		return gtk.ICON_SIZE_LARGE_TOOLBAR
	} else {
		return gtk.ICON_SIZE_SMALL_TOOLBAR
	}
}
