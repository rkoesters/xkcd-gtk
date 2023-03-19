package widget

import (
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type CacheWindow struct {
	*gtk.ApplicationWindow

	actions map[string]*glib.SimpleAction

	box                     *gtk.Box
	metadataLevelBar        *labeledLevelBar
	imageLevelBar           *labeledLevelBar
	downloadAllImagesButton *gtk.Button

	lastRefreshMetadata      time.Time
	lastRefreshMetadataMutex sync.RWMutex
	lastRefreshImages        time.Time
	lastRefreshImagesMutex   sync.RWMutex
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
	cw.box.SetMarginTop(style.PaddingAuxiliaryWindow)
	cw.box.SetMarginBottom(style.PaddingAuxiliaryWindow)
	cw.box.SetMarginStart(style.PaddingAuxiliaryWindow)
	cw.box.SetMarginEnd(style.PaddingAuxiliaryWindow)
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
	cw.imageLevelBar.SetMarginTop(style.PaddingAuxiliaryWindow)
	cw.box.PackStart(cw.imageLevelBar, false, true, 0)

	bb, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	bb.SetHAlign(gtk.ALIGN_END)
	bb.SetMarginTop(style.PaddingAuxiliaryWindow)
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
			cache.DownloadAllComicImages(func() cache.ViewRefresher { return cw })
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
	if cw == nil {
		return
	}
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

func (cw *CacheWindow) IsVisible() bool {
	if cw == nil {
		return false
	}
	return cw.ApplicationWindow.IsVisible()
}

const stalenessThreshold = 2 * time.Second

func (cw *CacheWindow) IsMetadataStale() bool {
	cw.lastRefreshMetadataMutex.RLock()
	defer cw.lastRefreshMetadataMutex.RUnlock()
	return time.Since(cw.lastRefreshMetadata) > stalenessThreshold
}

func (cw *CacheWindow) IsImagesStale() bool {
	cw.lastRefreshImagesMutex.RLock()
	defer cw.lastRefreshImagesMutex.RUnlock()
	return time.Since(cw.lastRefreshImages) > stalenessThreshold
}

func (cw *CacheWindow) RefreshMetadata() {
	if !cw.IsVisible() {
		return
	}
	if !cw.IsMetadataStale() {
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
	if !cw.IsVisible() {
		return
	}
	if !cw.IsMetadataStale() && !metadata.Complete() {
		return
	}

	cw.lastRefreshMetadataMutex.Lock()
	defer cw.lastRefreshMetadataMutex.Unlock()

	metaf, err := metadata.Fraction()
	if err != nil {
		log.Print("error refreshing cache window: ", err)
	}
	cw.metadataLevelBar.SetFraction(metaf)
	cw.metadataLevelBar.SetDetails(metadata.String())
	cw.lastRefreshMetadata = time.Now()
}

func (cw *CacheWindow) RefreshImages() {
	if !cw.IsVisible() {
		return
	}
	if !cw.IsImagesStale() {
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
	if !cw.IsVisible() {
		return
	}
	if !cw.IsImagesStale() && !images.Complete() {
		return
	}

	cw.lastRefreshImagesMutex.Lock()
	defer cw.lastRefreshImagesMutex.Unlock()

	imgf, err := images.Fraction()
	if err != nil {
		log.Print("error refreshing cache window: ", err)
	}
	cw.imageLevelBar.SetFraction(imgf)
	cw.imageLevelBar.SetDetails(images.String())
	cw.lastRefreshImages = time.Now()
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
	lb.SetMarginTop(4)
	lb.SetMarginBottom(4)
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
	log.Debugf("labeledLevelBar.SetFraction(%q)", f)
	lpb.bar.SetValue(f)
}

func (lpb *labeledLevelBar) SetDetails(s string) {
	log.Debugf("labeledLevelBar.SetDetails(%q)", s)
	lpb.details.SetText(s)
}
