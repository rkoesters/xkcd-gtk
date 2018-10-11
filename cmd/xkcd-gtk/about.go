package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

const (
	appID   = "com.github.rkoesters.xkcd-gtk"
	appName = "Comic Sticks"
)

var appVersion = "undefined"

var aboutDialog *gtk.AboutDialog

// ShowAboutDialog shows our application info to the user.
func (app *Application) ShowAboutDialog() {
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
		aboutDialog.SetComments("A simple xkcd viewer written in Go using GTK+3")
		aboutDialog.SetWebsite("https://github.com/rkoesters/xkcd-gtk")
		aboutDialog.SetCopyright("Copyright © 2015-2018 Ryan Koesters")
		aboutDialog.SetLicenseType(gtk.LICENSE_GPL_3_0)

		aboutDialog.SetAuthors([]string{"Ryan Koesters"})

		// We want to keep the about dialog around in case we want to
		// show it again.
		aboutDialog.HideOnDelete()
		aboutDialog.Connect("response", aboutDialog.Hide)
		aboutDialog.Connect("hide", func() {
			app.application.RemoveWindow(&aboutDialog.Window)
		})
	}

	// Set our parent window as the active window, but avoid
	// accidentally setting ourself as the parent window.
	win := app.application.GetActiveWindow()
	if win.Native() != aboutDialog.Native() {
		aboutDialog.SetTransientFor(win)
	}

	app.application.AddWindow(&aboutDialog.Window)
	aboutDialog.Present()
}
