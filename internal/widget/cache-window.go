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

	actions map[string]*glib.SimpleAction

	box                     *gtk.Box
	metadataLevelBar        *labeledLevelBar
	imageLevelBar           *labeledLevelBar
	downloadAllImagesButton *gtk.Button
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

	cw := &CacheWindow{
		ApplicationWindow: super,
		actions:           make(map[string]*glib.SimpleAction),
	}

	cw.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	cw.box.SetMarginTop(style.PaddingPropertiesDialog)
	cw.box.SetMarginBottom(style.PaddingPropertiesDialog)
	cw.box.SetMarginStart(style.PaddingPropertiesDialog)
	cw.box.SetMarginEnd(style.PaddingPropertiesDialog)
	cw.Add(cw.box)

	cw.metadataLevelBar, err = newLabeledLevelBar(l("Cached comic metadata"))
	if err != nil {
		return nil, err
	}
	cw.box.PackStart(cw.metadataLevelBar, false, true, 0)

	cw.imageLevelBar, err = newLabeledLevelBar(l("Cached comic images"))
	if err != nil {
		return nil, err
	}
	cw.box.PackStart(cw.imageLevelBar, false, true, 0)

	bb, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	bb.SetHAlign(gtk.ALIGN_END)
	cw.box.PackEnd(bb, false, true, 0)

	cw.downloadAllImagesButton, err = gtk.ButtonNewWithLabel(l("Download all comic images"))
	if err != nil {
		return nil, err
	}
	cw.downloadAllImagesButton.SetActionName("win.download-all-images")
	bb.PackStart(cw.downloadAllImagesButton, false, true, 0)

	registerAction := func(name string, fn interface{}) {
		action := glib.SimpleActionNew(name, nil)
		action.Connect("activate", fn)

		cw.actions[name] = action
		cw.AddAction(action)
	}

	registerAction("download-all-images", func() {
		cw.actions["download-all-images"].SetEnabled(false)
		go func() {
			cache.DownloadAllComicImages(cw)
			b := cw.actions["download-all-images"]
			glib.IdleAdd(func() {
				b.SetEnabled(true)
			})
		}()
	})

	cw.box.ShowAll()
	return cw, nil
}

func (cw *CacheWindow) Dispose() {
	cw.ApplicationWindow = nil
	cw.actions = nil
	cw.box = nil
	cw.metadataLevelBar = nil
	cw.imageLevelBar = nil
	cw.downloadAllImagesButton = nil
}

func (cw *CacheWindow) Present() {
	cw.ApplicationWindow.Present()
	go cw.RefreshMetadata()
	go cw.RefreshImages()
}

func (cw *CacheWindow) RefreshMetadata() {
	if !cw.IsVisible() {
		return
	}

	cs, err := cache.StatMetadata()
	if err != nil {
		log.Print("error refreshing cache window: ", err)
		return
	}

	glib.IdleAdd(func() {
		cw.RefreshMetadataWith(cs)
	})
}

func (cw *CacheWindow) RefreshMetadataWith(metadata cache.Stat) {
	metaf, err := metadata.Fraction()
	if err != nil {
		log.Print("error refreshing cache window: ", err)
	}
	cw.metadataLevelBar.SetFraction(metaf)
	cw.metadataLevelBar.SetDetails(metadata.String())
}

func (cw *CacheWindow) RefreshImages() {
	if !cw.IsVisible() {
		return
	}

	cs, err := cache.StatImages()
	if err != nil {
		log.Print("error refreshing cache window: ", err)
		return
	}

	glib.IdleAdd(func() {
		cw.RefreshImagesWith(cs)
	})
}

func (cw *CacheWindow) RefreshImagesWith(images cache.Stat) {
	imgf, err := images.Fraction()
	if err != nil {
		log.Print("error refreshing cache window: ", err)
	}
	cw.imageLevelBar.SetFraction(imgf)
	cw.imageLevelBar.SetDetails(images.String())
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
