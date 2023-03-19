package widget

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type CacheManager interface {
	MetadataStats() (int, int, error)
	ImageStats() (int, int, error)
	DownloadImages()
}

type CacheWindow struct {
	*gtk.ApplicationWindow
	box              *gtk.Box
	metadataLevelBar *labeledLevelBar
	imageLevelBar    *labeledLevelBar
}

var _ Widget = &CacheWindow{}

func NewCacheWindow(app Application) (*CacheWindow, error) {
	super, err := gtk.ApplicationWindowNew(app.GtkApplication())
	if err != nil {
		return nil, err
	}
	super.SetTitle(l("Cache manager"))
	super.SetSizeRequest(400, -1)
	super.HideOnDelete()
	super.Connect("hide", func(win gtk.IWindow) {
		app.RemoveWindow(win)
	})

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	box.SetMarginTop(style.PaddingPropertiesDialog)
	box.SetMarginBottom(style.PaddingPropertiesDialog)
	box.SetMarginStart(style.PaddingPropertiesDialog)
	box.SetMarginEnd(style.PaddingPropertiesDialog)
	super.Add(box)

	mpb, err := newLabeledLevelBar(l("Cached comic metadata"))
	if err != nil {
		return nil, err
	}
	box.PackStart(mpb, false, true, 0)
	ipb, err := newLabeledLevelBar(l("Cached comic images"))
	if err != nil {
		return nil, err
	}
	box.PackStart(ipb, false, true, 0)

	cw := &CacheWindow{
		ApplicationWindow: super,
		box:               box,
		metadataLevelBar:  mpb,
		imageLevelBar:     ipb,
	}
	cw.box.ShowAll()

	return cw, nil
}

func (cw *CacheWindow) Dispose() {
	cw.ApplicationWindow = nil
	cw.metadataLevelBar = nil
	cw.imageLevelBar = nil
}

func (cw *CacheWindow) Present() {
	cw.ApplicationWindow.Present()
	cw.RefreshMetadata()
	cw.RefreshImages()
}

func (cw *CacheWindow) RefreshMetadata() {
	go cw.refreshMetadata()
}

func (cw *CacheWindow) refreshMetadata() {
	if !cw.IsVisible() {
		return
	}

	cs, err := cache.StatMetadata()
	if err != nil {
		log.Print("error refreshing cache window: ", err)
		return
	}

	glib.IdleAdd(func() {
		metaf, err := cs.Fraction()
		if err != nil {
			log.Print("error refreshing cache window: ", err)
		}
		cw.metadataLevelBar.SetFraction(metaf)
		cw.metadataLevelBar.SetDetails(cs.String())
	})
}

func (cw *CacheWindow) RefreshImages() {
	go cw.refreshImages()
}

func (cw *CacheWindow) refreshImages() {
	if !cw.IsVisible() {
		return
	}

	cs, err := cache.StatImages()
	if err != nil {
		log.Print("error refreshing cache window: ", err)
		return
	}

	glib.IdleAdd(func() {
		imgf, err := cs.Fraction()
		if err != nil {
			log.Print("error refreshing cache window: ", err)
		}
		cw.imageLevelBar.SetFraction(imgf)
		cw.imageLevelBar.SetDetails(cs.String())
	})
}

type labeledLevelBar struct {
	*gtk.Box
	title   *gtk.Label
	bar     *gtk.LevelBar
	details *gtk.Label
}

func newLabeledLevelBar(title string) (*labeledLevelBar, error) {
	super, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}

	tl, err := gtk.LabelNew(title)
	if err != nil {
		return nil, err
	}
	tl.SetXAlign(0)
	tl.SetYAlign(1)
	super.PackStart(tl, false, true, 0)

	lb, err := gtk.LevelBarNew()
	if err != nil {
		return nil, err
	}
	lb.SetMode(gtk.LEVEL_BAR_MODE_CONTINUOUS)
	lb.RemoveOffsetValue(gtk.LEVEL_BAR_OFFSET_LOW)
	lb.RemoveOffsetValue(gtk.LEVEL_BAR_OFFSET_HIGH)
	lb.RemoveOffsetValue(gtk.LEVEL_BAR_OFFSET_FULL)
	lb.SetMarginTop(6)
	lb.SetMarginBottom(6)
	super.PackStart(lb, false, true, 0)

	dl, err := gtk.LabelNew("")
	if err != nil {
		return nil, err
	}
	dl.SetXAlign(1)
	dl.SetYAlign(0)
	super.PackStart(dl, false, true, 0)

	return &labeledLevelBar{
		Box:     super,
		title:   tl,
		bar:     lb,
		details: dl,
	}, nil
}

func (lpb *labeledLevelBar) SetFraction(f float64) {
	lpb.bar.SetValue(f)
}

func (lpb *labeledLevelBar) SetDetails(s string) {
	lpb.details.SetText(s)
}
