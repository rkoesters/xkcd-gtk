package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/build"
)

// NewAboutDialog creates our application's about dialog.
func NewAboutDialog(windowRemover func(gtk.IWindow)) (*gtk.AboutDialog, error) {
	dialog, err := gtk.AboutDialogNew()
	if err != nil {
		return nil, err
	}

	dialog.SetLogoIconName(build.AppID())
	dialog.SetProgramName(AppName())
	dialog.SetVersion(build.Version())
	dialog.SetComments(l("A simple xkcd viewer written in Go using GTK+3"))
	dialog.SetWebsite("https://github.com/rkoesters/xkcd-gtk")
	dialog.SetAuthors([]string{"Ryan Koesters"})
	dialog.SetTranslatorCredits(l("translator-credits"))
	dialog.SetCopyright("Copyright © 2015-2025 Ryan Koesters")
	dialog.SetLicenseType(gtk.LICENSE_GPL_3_0)

	// We want to keep the about dialog around in case we want to show it again,
	// so do not destroy it on close.
	dialog.HideOnDelete()
	dialog.Connect("response", func(dialog gtk.IWindow) {
		dialog.ToWindow().Hide()
	})
	dialog.Connect("hide", func(dialog gtk.IWindow) {
		windowRemover(dialog)
	})

	// Initialize our window accelerators.
	accels, err := gtk.AccelGroupNew()
	if err != nil {
		return nil, err
	}
	dialog.AddAccelGroup(accels)
	accels.Connect(gdk.KEY_w, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE, dialog.Close)

	return dialog, nil
}
