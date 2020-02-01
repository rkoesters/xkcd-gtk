package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
)

var l = glib.Local

func init() {
	glib.InitI18n(appID, paths.LocaleDir())
}
