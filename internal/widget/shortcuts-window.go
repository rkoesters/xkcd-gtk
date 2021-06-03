package widget

import (
	"github.com/gotk3/gotk3/gtk"
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
