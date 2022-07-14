package widget

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type ZoomBox struct {
	*gtk.ButtonBox

	zoomInButton    *gtk.Button
	zoomOutButton   *gtk.Button
	zoomResetButton *gtk.Button
}

var _ Widget = &ZoomBox{}

func NewZoomBox() (*ZoomBox, error) {
	const zbIconSize = gtk.ICON_SIZE_SMALL_TOOLBAR

	super, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	zb := &ZoomBox{
		ButtonBox: super,
	}

	zb.SetLayout(gtk.BUTTONBOX_EXPAND)
	zb.SetHomogeneous(false)

	zb.zoomOutButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	zb.zoomOutButton.SetTooltipText(l("Zoom out"))
	zb.zoomOutButton.SetActionName("win.zoom-out")
	zoomOutImg, err := gtk.ImageNewFromIconName("zoom-out-symbolic", zbIconSize)
	if err != nil {
		return nil, err
	}
	zb.zoomOutButton.SetImage(zoomOutImg)
	zb.PackStart(zb.zoomOutButton, true, true, 0)

	zb.zoomResetButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	zb.zoomResetButton.SetTooltipText(l("Reset zoom"))
	zb.zoomResetButton.SetActionName("win.zoom-reset")
	zb.PackStart(zb.zoomResetButton, true, true, 0)

	zb.zoomInButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	zb.zoomInButton.SetTooltipText(l("Zoom in"))
	zb.zoomInButton.SetActionName("win.zoom-in")
	zoomInImg, err := gtk.ImageNewFromIconName("zoom-in-symbolic", zbIconSize)
	if err != nil {
		return nil, err
	}
	zb.zoomInButton.SetImage(zoomInImg)
	zb.PackStart(zb.zoomInButton, true, true, 0)

	return zb, nil
}

func (zb *ZoomBox) Destroy() {
	if zb == nil {
		return
	}

	zb.ButtonBox = nil

	zb.zoomInButton = nil
	zb.zoomOutButton = nil
	zb.zoomResetButton = nil
}

func (zb *ZoomBox) IWidget() gtk.IWidget {
	return zb.ButtonBox
}

func (zb *ZoomBox) SetCurrentZoom(scale float64) error {
	zb.zoomResetButton.SetLabel(fmt.Sprintf("%.0f%%", scale*100))
	child, err := zb.zoomResetButton.GetChild()
	if err != nil {
		return err
	}
	label, err := gtk.WidgetToLabel(child.ToWidget())
	if err != nil {
		return err
	}
	label.SetXAlign(0.5) // align center
	// Zoom goes from `25%` to `500%`, so 4 characters. But
	// gtk_label_set_width_chars uses
	// pango_font_metrics_get_approximate_char_width which is not the widest
	// a character can be, so add 1 for padding to avoid resizing the widget
	// when changing the zoom level.
	label.SetWidthChars(5)
	return nil
}

func (zb *ZoomBox) SetCompact(compact bool) {
	if compact {
		zb.SetMarginStart(style.PaddingPopoverCompact)
		zb.SetMarginEnd(style.PaddingPopoverCompact)
	} else {
		zb.SetMarginStart(0)
		zb.SetMarginEnd(0)
	}
}
