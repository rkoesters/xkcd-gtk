// Package widget provides custom GTK+ widgets.
package widget

import (
	"github.com/gotk3/gotk3/gtk"
)

// Widget is a custom GTK+ widget.
type Widget interface {
	// IWidget returns the Widget's top-level gtk.Widget.
	IWidget() gtk.IWidget

	// Destroy performs clean up to aid garbage collection. Should
	// gracefully accept a nil receiver.
	Destroy()
}
