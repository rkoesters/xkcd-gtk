package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xdg/basedir"
	"github.com/skratchdot/open-golang/open"
	"log"
	"path/filepath"
)

const (
	appID   = "com.github.rkoesters.xkcd-gtk"
	appName = "XKCD Viewer"
)

// Application holds onto our GTK representation of our application.
type Application struct {
	GtkApp *gtk.Application
}

// NewApplication creates an instance of our GTK Application.
func NewApplication() (*Application, error) {
	var app Application
	var err error

	app.GtkApp, err = gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return nil, err
	}
	app.GtkApp.Connect("startup", app.LoadSearchIndex)
	app.GtkApp.Connect("activate", app.Activate)

	app.GtkApp.SetDefault()

	return &app, nil
}

// Activate creates and presents a new window to the user.
func (a *Application) Activate() {
	window, err := NewWindow(a)
	if err != nil {
		log.Fatal(err)
	}
	window.win.Present()
}

// CacheDir returns the path to our app's cache directory.
func CacheDir() string {
	return filepath.Join(basedir.CacheHome, appID)
}

// ConfigDir returns the path to our app's user configuration directory.
func ConfigDir() string {
	return filepath.Join(basedir.ConfigHome, appID)
}

// DataDir returns the path to our app's user data directory.
func DataDir() string {
	return filepath.Join(basedir.DataHome, appID)
}

// OpenURL is a simple function that can be connected to a GTK signal
// with an arbitrary URL to open that URL in the user's web browser.
func OpenURL(_ interface{}, url string) {
	err := open.Start(url)
	if err != nil {
		log.Print(err)
	}
}
