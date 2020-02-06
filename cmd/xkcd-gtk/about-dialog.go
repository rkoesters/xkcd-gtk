package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

// NewAboutDialog creates our about dialog.
func NewAboutDialog() (*gtk.AboutDialog, error) {
	dialog, err := gtk.AboutDialogNew()
	if err != nil {
		return nil, err
	}

	dialog.SetLogoIconName(appID)
	dialog.SetProgramName(appName)
	dialog.SetVersion(appVersion)
	dialog.SetComments(l("A simple xkcd viewer written in Go using GTK+3"))
	dialog.SetWebsite("https://github.com/rkoesters/xkcd-gtk")
	dialog.SetAuthors([]string{"Ryan Koesters"})
	dialog.SetCopyright("Copyright © 2015-2020 Ryan Koesters")
	dialog.SetLicenseType(gtk.LICENSE_GPL_3_0)

	return dialog, nil
}

var aboutDialog *gtk.AboutDialog

// ShowAbout shows our application info to the user.
func (app *Application) ShowAbout() {
	var err error

	if aboutDialog == nil {
		aboutDialog, err = NewAboutDialog()
		if err != nil {
			log.Print("error creating about dialog: ", err)
			return
		}

		// We want to keep the about dialog around in case we want to
		// show it again.
		aboutDialog.HideOnDelete()
		aboutDialog.Connect("response", aboutDialog.Hide)
		aboutDialog.Connect("hide", func() {
			app.application.RemoveWindow(&aboutDialog.Window)
		})
	}

	// Set our parent window as the active window, but avoid accidentally
	// setting ourself as the parent window.
	win := app.application.GetActiveWindow()
	if win.Native() != aboutDialog.Native() {
		aboutDialog.SetTransientFor(win)
	}

	app.application.AddWindow(&aboutDialog.Window)
	aboutDialog.Present()
}
