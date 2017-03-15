package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xdg/basedir"
	"log"
	"path/filepath"
)

const appId = "com.ryankoesters.xkcd-gtk"

type Application struct {
	GtkApp *gtk.Application
}

func NewApplication() (*Application, error) {
	var app Application
	var err error

	app.GtkApp, err = gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return nil, err
	}
	app.GtkApp.Connect("activate", app.Activate)

	glib.SetApplicationName("XKCD Viewer")
	gtk.WindowSetDefaultIconName("xkcd-gtk")
	app.GtkApp.SetDefault()

	return &app, nil
}

func (a *Application) Activate() {
	window, err := NewWindow(a)
	if err != nil {
		log.Fatal(err)
	}
	window.win.Present()
}

func CacheDir() string {
	return filepath.Join(basedir.CacheHome, appId)
}

func ConfigDir() string {
	return filepath.Join(basedir.ConfigHome, appId)
}

func DataDir() string {
	return filepath.Join(basedir.DataHome, appId)
}
