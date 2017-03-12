package main

import (
	"github.com/gotk3/gotk3/gtk"
)

// NewAboutDialog creates a gtk.AboutDialog with our app's info.
func NewAboutDialog(parent *Window) (*gtk.AboutDialog, error) {
	abt, err := gtk.AboutDialogNew()
	if err != nil {
		return nil, err
	}
	abt.SetTransientFor(parent.win)

	abt.SetLogoIconName("xkcd-gtk")
	abt.SetProgramName("XKCD Viewer")
	abt.SetVersion("0.2")
	abt.SetComments("A simple XKCD comic reader for GNOME")
	abt.SetWebsite("https://github.com/rkoesters/xkcd-gtk")
	abt.SetCopyright("Copyright © 2015-2017 Ryan Koesters")
	abt.SetLicenseType(gtk.LICENSE_GPL_3_0)

	abt.SetAuthors([]string{"Ryan Koesters"})

	return abt, nil
}
