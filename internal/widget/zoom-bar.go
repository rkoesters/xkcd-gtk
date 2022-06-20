package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type ZoomBar struct {
	accels *gtk.AccelGroup // ptr to win.accels

	box *gtk.ButtonBox

	zoomInButton  *gtk.Button
	zoomOutButton *gtk.Button

	comicContainer *ImageViewer
}

var _ Widget = &ZoomBar{}

func NewZoomBar(accels *gtk.AccelGroup, comicContainer *ImageViewer) (*ZoomBar, error) {
	var err error

	zb := &ZoomBar{
		accels:         accels,
		comicContainer: comicContainer,
	}

	zb.box, err = gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	zb.box.SetLayout(gtk.BUTTONBOX_EXPAND)

	zb.zoomOutButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	zb.zoomOutButton.SetTooltipText(l("Zoom out"))
	zb.zoomOutButton.SetProperty("action-name", "win.zoom-out")
	zb.zoomOutButton.AddAccelerator("activate", zb.accels, gdk.KEY_minus, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
	zb.box.Add(zb.zoomOutButton)

	zb.zoomInButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	zb.zoomInButton.SetTooltipText(l("Zoom in"))
	zb.zoomInButton.SetProperty("action-name", "win.zoom-in")
	zb.zoomInButton.AddAccelerator("activate", zb.accels, gdk.KEY_equal, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
	zb.box.Add(zb.zoomInButton)

	return zb, nil
}

func (zb *ZoomBar) Destroy() {
	zb.accels = nil

	zb.box = nil

	zb.zoomOutButton = nil
	zb.zoomInButton = nil

	zb.comicContainer = nil
}

func (zb *ZoomBar) IWidget() gtk.IWidget {
	return zb.box
}

func (zb *ZoomBar) SetZoomInButtonImage(image *gtk.Image) {
	zb.zoomInButton.SetImage(image)
}

func (nb *ZoomBar) SetZoomOutButtonImage(image *gtk.Image) {
	nb.zoomOutButton.SetImage(image)
}

func (nb *ZoomBar) SetLinkedButtons(linked bool) {
	sc, err := nb.box.GetStyleContext()
	if err != nil {
		log.Print(err)
		return
	}

	if linked {
		sc.AddClass(style.ClassLinked)
		nb.box.SetSpacing(0)
	} else {
		sc.RemoveClass(style.ClassLinked)
		nb.box.SetSpacing(4)
	}
}
