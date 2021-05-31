package widget

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type ImageViewer struct {
	ScrolledWindow    *gtk.ScrolledWindow
	ScrolledWindowCtx *gtk.StyleContext
	Image             *gtk.Image
}

func NewImageViewer() (*ImageViewer, error) {
	var err error

	iv := new(ImageViewer)

	iv.ScrolledWindow, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}

	iv.ScrolledWindow.SetSizeRequest(500, 400)

	iv.ScrolledWindowCtx, err = iv.ScrolledWindow.GetStyleContext()
	if err != nil {
		return nil, err
	}
	iv.ScrolledWindowCtx.AddClass(style.ClassComicContainer)

	iv.Image, err = gtk.ImageNew()
	if err != nil {
		return nil, err
	}
	iv.Image.SetHAlign(gtk.ALIGN_CENTER)
	iv.Image.SetVAlign(gtk.ALIGN_CENTER)

	iv.ScrolledWindow.Add(iv.Image)

	return iv, nil
}

func (iv *ImageViewer) Destroy() {
	iv.ScrolledWindow = nil
	iv.ScrolledWindowCtx = nil
	iv.Image = nil
}

func (iv *ImageViewer) Show() {
	iv.ScrolledWindow.ShowAll()
}
