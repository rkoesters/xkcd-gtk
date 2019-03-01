package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

// NewShortcutsWindow creates a gtk.ShortcutsWindow populated with our
// application's keyboard shortcuts.
func NewShortcutsWindow() (*gtk.ShortcutsWindow, error) {
	builder, err := gtk.BuilderNew()
	if err != nil {
		return nil, err
	}

	err = builder.AddFromString(shortcutsWindowUI)
	if err != nil {
		return nil, err
	}

	obj, err := builder.GetObject("shortcuts-window")
	if err != nil {
		return nil, err
	}

	return obj.(*gtk.ShortcutsWindow), nil
}

var shortcutsWindow *gtk.ShortcutsWindow

// ShowShortcuts shows a shortcuts window to the user.
func (app *Application) ShowShortcuts() {
	var err error
	if shortcutsWindow == nil {
		shortcutsWindow, err = NewShortcutsWindow()
		if err != nil {
			log.Print(err)
			return
		}

		// We want to keep the shortcuts window around in case
		// we want to show it again.
		shortcutsWindow.HideOnDelete()
		shortcutsWindow.Connect("hide", func() {
			app.application.RemoveWindow(&shortcutsWindow.Window)
		})
	}

	app.application.AddWindow(&shortcutsWindow.Window)
	shortcutsWindow.Present()
}
