package widget

import (
	"github.com/gotk3/gotk3/gtk"
)

type Window interface {
	Widget

	// IWindow returns the Window's underlying gtk.Window or
	// gtk.ApplicationWindow.
	IWindow() gtk.IWindow

	// Close requests that the window be closed.
	Close()
}
