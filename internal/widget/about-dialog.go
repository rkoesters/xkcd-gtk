package widget

import (
	"github.com/gotk3/gotk3/gtk"
)

// NewAboutDialog creates our about dialog.
func NewAboutDialog(icon, name, version string) (*gtk.AboutDialog, error) {
	dialog, err := gtk.AboutDialogNew()
	if err != nil {
		return nil, err
	}

	dialog.SetLogoIconName(icon)
	dialog.SetProgramName(name)
	dialog.SetVersion(version)
	dialog.SetComments(l("A simple xkcd viewer written in Go using GTK+3"))
	dialog.SetWebsite("https://github.com/rkoesters/xkcd-gtk")
	dialog.SetAuthors([]string{"Ryan Koesters"})
	dialog.SetCopyright("Copyright © 2015-2022 Ryan Koesters")
	dialog.SetLicenseType(gtk.LICENSE_GPL_3_0)

	return dialog, nil
}
