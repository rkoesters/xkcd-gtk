package main

import (
	"github.com/gotk3/gotk3/glib"
)

var l = glib.Local

func init() {
	glib.InitI18n(appID, LocaleDir())
}
