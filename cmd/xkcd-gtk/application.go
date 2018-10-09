package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xdg/basedir"
	"log"
	"path/filepath"
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
	app.GtkApp.Connect("startup", app.LoadCSS)
	app.GtkApp.Connect("startup", app.LoadSearchIndex)
	app.GtkApp.Connect("activate", app.Activate)

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
