// Package widget provides custom GTK+ widgets.
package widget

import (
	"github.com/gotk3/gotk3/gtk"
)

// Widget is a custom GTK+ widget.
type Widget interface {
	// Our custom widgets should embed a gtk.Widget.
	gtk.IWidget

	// Dispose performs clean up to aid garbage collection. Should break
	// reference cycles, if any. Must gracefully accept a nil receiver.
	Dispose()
}
