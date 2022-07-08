package widget

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type ZoomBox struct {
	box *gtk.ButtonBox

	zoomInButton    *gtk.Button
	zoomOutButton   *gtk.Button
	zoomResetButton *gtk.Button
}

var _ Widget = &ZoomBox{}

func NewZoomBox() (*ZoomBox, error) {
	const zbIconSize = gtk.ICON_SIZE_SMALL_TOOLBAR

	var err error

	zb := &ZoomBox{}

	zb.box, err = gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	zb.box.SetLayout(gtk.BUTTONBOX_EXPAND)
	zb.box.SetProperty("homogeneous", false)

	zb.zoomOutButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	zb.zoomOutButton.SetTooltipText(l("Zoom out"))
	zb.zoomOutButton.SetProperty("action-name", "win.zoom-out")
	zoomOutImg, err := gtk.ImageNewFromIconName("zoom-out-symbolic", zbIconSize)
	if err != nil {
		return nil, err
	}
	zb.zoomOutButton.SetImage(zoomOutImg)
	zb.box.PackStart(zb.zoomOutButton, true, true, 0)

	zb.zoomResetButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	zb.zoomResetButton.SetTooltipText(l("Reset zoom"))
	zb.zoomResetButton.SetProperty("action-name", "win.zoom-reset")
	zb.box.PackStart(zb.zoomResetButton, true, true, 0)

	zb.zoomInButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	zb.zoomInButton.SetTooltipText(l("Zoom in"))
	zb.zoomInButton.SetProperty("action-name", "win.zoom-in")
	zoomInImg, err := gtk.ImageNewFromIconName("zoom-in-symbolic", zbIconSize)
	if err != nil {
		return nil, err
	}
	zb.zoomInButton.SetImage(zoomInImg)
	zb.box.PackStart(zb.zoomInButton, true, true, 0)

	return zb, nil
}

func (zb *ZoomBox) Destroy() {
	zb.box = nil

	zb.zoomInButton = nil
	zb.zoomOutButton = nil
	zb.zoomResetButton = nil
}

func (zb *ZoomBox) IWidget() gtk.IWidget {
	return zb.box
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
	// Zoom goes from `25%` to `500%`, so 4 characters.
	label.SetWidthChars(4)
	return nil
}

func (zb *ZoomBox) SetCompact(compact bool) {
	if compact {
		zb.box.SetMarginStart(style.PaddingPopoverCompact)
		zb.box.SetMarginEnd(style.PaddingPopoverCompact)
	} else {
		zb.box.SetMarginStart(0)
		zb.box.SetMarginEnd(0)
	}
}
