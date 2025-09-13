package widget

import (
	"errors"
	"math"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/state"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type ImageViewer struct {
	*gtk.ScrolledWindow

	image          *gtk.Image
	unscaledPixbuf *gdk.Pixbuf // will be inverted in dark mode
	scale          float64
	finalPixbuf    *gdk.Pixbuf // displayed to the user

	eventBox *gtk.EventBox

	contextMenu *ContextMenu
}

var _ Widget = &ImageViewer{}

func NewImageViewer(actionGroup glib.IActionGroup, imageScale float64, bookmarkedGetter func() bool, bookmarkedSetter func(bool)) (*ImageViewer, error) {
	super, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	iv := &ImageViewer{
		ScrolledWindow: super,

		scale: safeScale(imageScale),
	}

	iv.SetSizeRequest(500, 400)

	sc, err := iv.GetStyleContext()
	if err != nil {
		return nil, err
	}
	sc.AddClass(style.ClassComicContainer)

	iv.image, err = gtk.ImageNew()
	if err != nil {
		return nil, err
	}
	iv.image.SetHAlign(gtk.ALIGN_CENTER)
	iv.image.SetVAlign(gtk.ALIGN_CENTER)

	iv.eventBox, err = gtk.EventBoxNew()
	if err != nil {
		return nil, err
	}
	iv.eventBox.Add(iv.image)
	iv.Add(iv.eventBox)

	iv.contextMenu, err = NewContextMenu(iv.eventBox, actionGroup, bookmarkedGetter, bookmarkedSetter)
	if err != nil {
		return nil, err
	}

	iv.eventBox.Connect("button-press-event", func(eventBox *gtk.EventBox, event *gdk.Event) bool {
		button := gdk.EventButtonNewFromEvent(event)
		switch button.Button() {
		case gdk.BUTTON_SECONDARY:
			iv.contextMenu.PopupAtPointer(button)
			return true
		default:
			return false
		}
	})

	iv.ShowAll()

	return iv, nil
}

func (iv *ImageViewer) Dispose() {
	if iv == nil {
		return
	}

	iv.ScrolledWindow = nil

	iv.image = nil
	iv.unscaledPixbuf = nil
	iv.finalPixbuf = nil
	iv.eventBox = nil

	iv.contextMenu.Dispose()
	iv.contextMenu = nil
}

func (iv *ImageViewer) ShowLoadingScreen() {
	iv.image.SetFromIconName("image-loading-symbolic", gtk.ICON_SIZE_DIALOG)
}

func (iv *ImageViewer) SetScale(scale float64) float64 {
	var err error
	iv.scale = safeScale(scale)
	iv.finalPixbuf, err = scaleImage(iv.unscaledPixbuf, iv.scale)
	if err != nil {
		log.Print(err)
		return 0
	}
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

func (iv *ImageViewer) DrawComic(comicId int, darkMode bool) error {
	log.Debugf("DrawComic(id=%v, darkMode=%v)", comicId, darkMode)
	path := cache.ComicImagePath(comicId)
	var err error
	iv.unscaledPixbuf, err = gdk.PixbufNewFromFile(path)
	if err != nil {
		return err
	}
	if darkMode {
		err = iv.applyDarkModeImageInversion()
		if err != nil {
			return err
		}
	}
	iv.finalPixbuf, err = scaleImage(iv.unscaledPixbuf, iv.scale)
	if err != nil {
		return err
	}
	iv.image.SetFromPixbuf(iv.finalPixbuf)
	return nil
}

func (iv *ImageViewer) applyDarkModeImageInversion() error {
	pixels := iv.unscaledPixbuf.GetPixels()
	colorspace := iv.unscaledPixbuf.GetColorspace()
	alpha := iv.unscaledPixbuf.GetHasAlpha()
	bitsPerSample := iv.unscaledPixbuf.GetBitsPerSample()
	width := iv.unscaledPixbuf.GetWidth()
	height := iv.unscaledPixbuf.GetHeight()
	rowstride := iv.unscaledPixbuf.GetRowstride()
	nChannels := iv.unscaledPixbuf.GetNChannels()
	log.Debugf("inverting comic image: len(pixels) = %v, colorspace = %v, alpha = %v, bitsPerSample = %v, width = %v, height = %v, rowstride = %v, nChannels = %v", len(pixels), colorspace, alpha, bitsPerSample, width, height, rowstride, nChannels)
	
	for y := 0; y < height; y++ {
		rowstart := y * rowstride
		for x := 0; x < width; x++ {
			index := rowstart + (x * nChannels)
			switch nChannels {
			case 3, 4:
				h, s, l := rgbToHsl(pixels[index], pixels[index+1], pixels[index+2])
				l = 1 - l
				r, g, b := hslToRgb(h, s, l)
				pixels[index] = r
				pixels[index+1] = g
				pixels[index+2] = b
			default:
				return errors.New("unsupported number of channels")
			}
		}
	}
	return nil
}

func (iv *ImageViewer) SetTooltipText(s string) {
	iv.image.SetTooltipText(s)
}

func scaleImage(unscaled *gdk.Pixbuf, scale float64) (*gdk.Pixbuf, error) {
	width := int(float64(unscaled.GetWidth()) * scale)
	height := int(float64(unscaled.GetHeight()) * scale)
	return unscaled.ScaleSimple(width, height, gdk.INTERP_BILINEAR)
}

func safeScale(scale float64) float64 {
	switch {
	case scale < state.ImageScaleMin:
		return state.ImageScaleMin
	case scale > state.ImageScaleMax:
		return state.ImageScaleMax
	default:
		return scale
	}
}

func rgbToHsl(r, g, b uint8) (h, s, l float64) {
	rf := float64(r) / math.MaxUint8
	gf := float64(g) / math.MaxUint8
	bf := float64(b) / math.MaxUint8

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	c := max - min
	l = max - (c / 2)

	if l == 0 || l == 1 {
		s = 0
	} else {
		s = (max - l) / math.Min(l, 1 - l)
	}

	if max == min {
		h = 0
	} else if max == rf {
		h = 60 * math.Mod((gf - bf) / c, 6)
	} else if max == gf {
		h = 60 * ((bf - rf) / c + 2)
	} else {
		h = 60 * ((rf - gf) / c + 4)
	}

	if h < 0 {
		h += 360
	}
	return
}

func hslToRgb(h, s, l float64) (r, g, b uint8) {
	a := s * math.Min(l, 1 - l)
	f := func(n float64) uint8 {
		k := math.Mod((n + (h / 30)), 12)
		fn := l - a * math.Max(-1, math.Min(k - 3, math.Min(9 - k, 1)))
		return uint8(math.Round(fn * math.MaxUint8))
	}

	r = f(0)
	g = f(8)
	b = f(4)
	return
}
