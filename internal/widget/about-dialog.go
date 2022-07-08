package widget

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/build"
)

// NewAboutDialog creates our about dialog.
func NewAboutDialog(windowRemover func(gtk.IWindow)) (*gtk.AboutDialog, error) {
	dialog, err := gtk.AboutDialogNew()
	if err != nil {
		return nil, err
	}

	dialog.SetLogoIconName(build.AppID)
	dialog.SetProgramName(AppName())
	dialog.SetVersion(build.Version())
	dialog.SetComments(l("A simple xkcd viewer written in Go using GTK+3"))
	dialog.SetWebsite("https://github.com/rkoesters/xkcd-gtk")
	dialog.SetAuthors([]string{"Ryan Koesters"})
	dialog.SetCopyright("Copyright © 2015-2022 Ryan Koesters")
	dialog.SetLicenseType(gtk.LICENSE_GPL_3_0)

	// We want to keep the about dialog around in case we want to show it
	// again.
	dialog.HideOnDelete()
	dialog.Connect("response", dialog.Hide)
	dialog.Connect("hide", func() {
		windowRemover(&dialog.Window)
	})

	return dialog, nil
}
