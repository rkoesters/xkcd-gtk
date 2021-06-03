// Package widget provides custom GTK+ widgets.
package widget

import (
	"github.com/gotk3/gotk3/gtk"
)

type Widget interface {
	// Return the Widget's top-level gtk.Widget.
	IWidget() gtk.IWidget
}
