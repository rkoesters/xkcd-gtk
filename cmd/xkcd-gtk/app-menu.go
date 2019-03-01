package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

// NewAppMenu creates a glib.MenuModel populated with our application's
// app menu.
func NewAppMenu() (*glib.MenuModel, error) {
	builder, err := gtk.BuilderNew()
	if err != nil {
		return nil, err
	}

	err = builder.AddFromString(appMenuUI)
	if err != nil {
		return nil, err
	}

	obj, err := builder.GetObject("app-menu")
	if err != nil {
		return nil, err
	}

	return obj.(*glib.MenuModel), nil
}

// SetupAppMenu creates an AppMenu if the environment wants it.
func (app *Application) SetupAppMenu() {
	if app.application.PrefersAppMenu() {
		menu, err := NewAppMenu()
		if err != nil {
			log.Fatal(err)
		}
		app.application.SetAppMenu(menu)
	}
}
