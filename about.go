package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

const (
	appID      = "com.github.rkoesters.xkcd-gtk"
	appName    = "Comic Sticks"
	appVersion = "0.9.9"
)

var aboutDialog *gtk.AboutDialog

// ShowAboutDialog shows our application info to the user.
func (a *Application) ShowAboutDialog() {
	var err error
	if aboutDialog == nil {
		aboutDialog, err = gtk.AboutDialogNew()
		if err != nil {
			log.Print(err)
			return
		}

		aboutDialog.SetLogoIconName(appID)
		aboutDialog.SetProgramName(appName)
		aboutDialog.SetVersion(appVersion)
		aboutDialog.SetComments("A simple xkcd viewer written in Go using GTK3")
		aboutDialog.SetWebsite("https://github.com/rkoesters/xkcd-gtk")
		aboutDialog.SetCopyright("Copyright © 2015-2018 Ryan Koesters")
		aboutDialog.SetLicenseType(gtk.LICENSE_GPL_3_0)

		aboutDialog.SetAuthors([]string{"Ryan Koesters"})

		// We want to keep the about dialog around in case we want to
		// show it again.
		aboutDialog.HideOnDelete()
		aboutDialog.Connect("response", aboutDialog.Hide)
		aboutDialog.Connect("hide", func() {
			a.GtkApp.RemoveWindow(&aboutDialog.Window)
		})
	}
	a.GtkApp.AddWindow(&aboutDialog.Window)
	aboutDialog.Present()
}
