package widget

import (
	"errors"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// NewShortcutsWindow creates a gtk.ShortcutsWindow populated with our
// application's keyboard shortcuts.
func NewShortcutsWindow(windowRemover func(gtk.IWindow)) (*gtk.ShortcutsWindow, error) {
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

	sw, ok := obj.(*gtk.ShortcutsWindow)
	if !ok {
		return nil, errors.New("error converting shortcuts-window into *gtk.ShortcutsWindow")
	}

	// We want to keep the shortcuts window around in case we want to show
	// it again, so do not destroy it on close.
	sw.HideOnDelete()
	sw.Connect("hide", func() {
		windowRemover(&sw.Window)
	})

	// Initialize our window accelerators.
	accels, err := gtk.AccelGroupNew()
	if err != nil {
		return nil, err
	}
	sw.AddAccelGroup(accels)
	accels.Connect(gdk.KEY_w, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE, sw.Close)

	return sw, nil
}
