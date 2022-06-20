package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/style"
	"math"
)

type ImageViewer struct {
	scrolledWindow    *gtk.ScrolledWindow
	scrolledWindowCtx *gtk.StyleContext

	image          *gtk.Image
	unscaledPixbuf *gdk.Pixbuf // will be inverted in dark mode
	scale          float64
	finalPixbuf    *gdk.Pixbuf // displayed to the user

	eventBox *gtk.EventBox

	contextMenu *ContextMenu
}

var _ Widget = &ImageViewer{}

func NewImageViewer(parent *gtk.ApplicationWindow, imageScale float64) (*ImageViewer, error) {
	var err error

	iv := new(ImageViewer)

	iv.scale = safeScale(imageScale)

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

	iv.contextMenu, err = NewContextMenu(parent.IActionGroup)
	if err != nil {
		return nil, err
	}

	iv.eventBox, err = gtk.EventBoxNew()
	if err != nil {
		return nil, err
	}
	iv.eventBox.Add(iv.image)
	iv.eventBox.Connect("button-press-event", func(eventBox *gtk.EventBox, event *gdk.Event) bool {
		button := gdk.EventButtonNewFromEvent(event)
		switch button.Button() {
		case gdk.BUTTON_SECONDARY:
			iv.contextMenu.Present(event)
			return true
		default:
			return false
		}
	})
	iv.scrolledWindow.Add(iv.eventBox)

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
	iv.unscaledPixbuf = nil
	iv.finalPixbuf = nil
	iv.eventBox = nil

	iv.contextMenu.Destroy()
	iv.contextMenu = nil
}

func (iv *ImageViewer) Show() {
	iv.scrolledWindow.ShowAll()
}

func (iv *ImageViewer) ShowLoadingScreen() {
	iv.image.SetFromIconName("image-loading-symbolic", gtk.ICON_SIZE_DIALOG)
}

func (iv *ImageViewer) SetScale(scale float64) float64 {
	iv.scale = safeScale(scale)
	iv.applyImageScaling()
	iv.image.SetFromPixbuf(iv.finalPixbuf)
	return iv.scale
}

const zoomIncrement = 0.25

func (iv *ImageViewer) ZoomIn() float64 {
	return iv.SetScale(iv.scale + zoomIncrement)
}

func (iv *ImageViewer) ZoomOut() float64 {
	return iv.SetScale(iv.scale - zoomIncrement)
}

func (iv *ImageViewer) SetComic(comicId int, darkMode bool) {
	path := cache.ComicImagePath(comicId)
	var err error
	iv.unscaledPixbuf, err = gdk.PixbufNewFromFile(path)
	if err != nil {
		log.Print(err)
		return
	}
	iv.applyDarkModeImageInversion(darkMode)
	err = iv.applyImageScaling()
	if err != nil {
		log.Print(err)
		return
	}
	iv.image.SetFromPixbuf(iv.finalPixbuf)
}

func (iv *ImageViewer) applyDarkModeImageInversion(enabled bool) {
	if enabled {
		// Invert the pixels of the comic image.
		pixels := iv.unscaledPixbuf.GetPixels()
		for i := 0; i < len(pixels); i++ {
			pixels[i] = math.MaxUint8 - pixels[i]
		}
	}
}

func (iv *ImageViewer) applyImageScaling() error {
	var unscaledWidth, unscaledHeight, width, height int
	unscaledWidth = iv.unscaledPixbuf.GetWidth()
	unscaledHeight = iv.unscaledPixbuf.GetHeight()
	width = int(float64(unscaledWidth) * iv.scale)
	height = int(float64(unscaledHeight) * iv.scale)
	var err error
	iv.finalPixbuf, err = iv.unscaledPixbuf.ScaleSimple(width, height, gdk.INTERP_HYPER)
	return err
}

func (iv *ImageViewer) SetTooltipText(s string) {
	iv.image.SetTooltipText(s)
}

func safeScale(scale float64) float64 {
	switch {
	case scale < ImageScaleMin:
		return ImageScaleMin
	case scale > ImageScaleMax:
		return ImageScaleMax
	default:
		return scale
	}
}
