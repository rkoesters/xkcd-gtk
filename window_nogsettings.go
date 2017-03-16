// +build windows darwin

package main

import (
	"github.com/gotk3/gotk3/gtk"
)

func lookupHeaderBarIconSize() gtk.IconSize {
	return gtk.ICON_SIZE_LARGE_TOOLBAR
}
