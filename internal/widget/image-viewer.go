package widget

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
	"log"
	"math"
)

type ImageViewer struct {
	scrolledWindow    *gtk.ScrolledWindow
	scrolledWindowCtx *gtk.StyleContext
	image             *gtk.Image
}

var _ Widget = &ImageViewer{}

func NewImageViewer() (*ImageViewer, error) {
	var err error

	iv := new(ImageViewer)

	iv.scrolledWindow, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}

	iv.scrolledWindow.SetSizeRequest(500, 400)

	iv.scrolledWindowCtx, err = iv.scrolledWindow.GetStyleContext()
	if err != nil {
		return nil, err
	}
	iv.scrolledWindowCtx.AddClass(style.ClassComicContainer)

	iv.image, err = gtk.ImageNew()
	if err != nil {
		return nil, err
	}
	iv.image.SetHAlign(gtk.ALIGN_CENTER)
	iv.image.SetVAlign(gtk.ALIGN_CENTER)

	iv.scrolledWindow.Add(iv.image)

	return iv, nil
}

func (iv *ImageViewer) IWidget() gtk.IWidget {
	// Return the top-level widget.
	return iv.scrolledWindow
}

func (iv *ImageViewer) Destroy() {
	iv.scrolledWindow = nil
	iv.scrolledWindowCtx = nil
	iv.image = nil
}

func (iv *ImageViewer) Show() {
	iv.scrolledWindow.ShowAll()
}

func (iv *ImageViewer) ShowLoadingScreen() {
	iv.image.SetFromIconName("image-loading-symbolic", gtk.ICON_SIZE_DIALOG)
}

func (iv *ImageViewer) SetFromFile(path string, darkMode bool) {
	iv.image.SetFromFile(path)
	iv.applyDarkModeImageInversion(darkMode)
}

func (iv *ImageViewer) applyDarkModeImageInversion(enabled bool) {
	if enabled {
		// Invert the pixels of the comic image.
		pixbuf := iv.image.GetPixbuf()
		if pixbuf == nil {
			log.Print("pixbuf == nil")
			return
		}
		pixels := pixbuf.GetPixels()
		for i := 0; i < len(pixels); i++ {
			pixels[i] = math.MaxUint8 - pixels[i]
		}
	}
}

func (iv *ImageViewer) SetTooltipText(s string) {
	iv.image.SetTooltipText(s)
}
