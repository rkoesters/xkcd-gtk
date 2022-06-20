package widget

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/log"
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
	zb.box.Add(zb.zoomOutButton)

	zb.zoomResetButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	zb.zoomResetButton.SetTooltipText(l("Reset zoom"))
	zb.zoomResetButton.SetProperty("action-name", "win.zoom-reset")
	zb.box.Add(zb.zoomResetButton)

	zb.zoomInButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	zb.zoomInButton.SetTooltipText(l("Zoom in"))
	zb.zoomInButton.SetProperty("action-name", "win.zoom-in")
	zb.box.Add(zb.zoomInButton)

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

func (zb *ZoomBox) SetZoomInButtonImage(image *gtk.Image) {
	zb.zoomInButton.SetImage(image)
}

func (zb *ZoomBox) SetZoomOutButtonImage(image *gtk.Image) {
	zb.zoomOutButton.SetImage(image)
}

func (zb *ZoomBox) SetCurrentZoom(scale float64) {
	zb.zoomResetButton.SetLabel(fmt.Sprintf("%.0f%%", scale*100))
}

func (zb *ZoomBox) SetLinkedButtons(linked bool) {
	sc, err := zb.box.GetStyleContext()
	if err != nil {
		log.Print(err)
		return
	}

	if linked {
		sc.AddClass(style.ClassLinked)
		zb.box.SetSpacing(0)
	} else {
		sc.RemoveClass(style.ClassLinked)
		zb.box.SetSpacing(4)
	}
}
