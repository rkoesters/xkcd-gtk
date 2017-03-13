package main

import (
	"github.com/gotk3/gotk3/gtk"
)

// NewAboutDialog creates a gtk.AboutDialog with our app's info.
func NewAboutDialog() (*gtk.AboutDialog, error) {
	abt, err := gtk.AboutDialogNew()
	if err != nil {
		return nil, err
	}

	abt.SetLogoIconName("xkcd-gtk")
	abt.SetProgramName("XKCD Viewer")
	abt.SetVersion("0.3.1")
	abt.SetComments("A simple XKCD comic reader for GNOME")
	abt.SetWebsite("https://github.com/rkoesters/xkcd-gtk")
	abt.SetCopyright("Copyright © 2015-2017 Ryan Koesters")
	abt.SetLicenseType(gtk.LICENSE_GPL_3_0)

	abt.SetAuthors([]string{"Ryan Koesters"})

	abt.Connect("response", aboutDialogClose, abt)

	return abt, nil
}

func aboutDialogClose(abt *gtk.AboutDialog) {
	abt.Close()
}
