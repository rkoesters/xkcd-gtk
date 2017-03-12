package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd"
	"log"
	"math/rand"
	"strconv"
)

// Window is the main application window.
type Window struct {
	comic    *xkcd.Comic
	win      *gtk.ApplicationWindow
	hdr      *gtk.HeaderBar
	previous *gtk.Button
	next     *gtk.Button
	img      *gtk.Image
}

// New creates a new XKCD viewer window.
func NewWindow(app *Application) (*Window, error) {
	var err error

	w := new(Window)

	w.win, err = gtk.ApplicationWindowNew(app.GtkApp)
	if err != nil {
		return nil, err
	}
	w.win.SetDefaultSize(1000, 800)

	// Create HeaderBar
	w.hdr, err = gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}
	w.hdr.SetTitle("XKCD Viewer")
	w.hdr.SetShowCloseButton(true)

	navBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}
	navBoxStyleContext, err := navBox.GetStyleContext()
	if err != nil {
		return nil, err
	}
	navBoxStyleContext.AddClass("linked")

	w.previous, err = gtk.ButtonNewFromIconName("go-previous-symbolic", gtk.ICON_SIZE_LARGE_TOOLBAR)
	if err != nil {
		return nil, err
	}
	w.previous.Connect("clicked", w.PreviousComic)
	navBox.Add(w.previous)

	w.next, err = gtk.ButtonNewFromIconName("go-next-symbolic", gtk.ICON_SIZE_LARGE_TOOLBAR)
	if err != nil {
		return nil, err
	}
	w.next.Connect("clicked", w.NextComic)
	navBox.Add(w.next)

	w.hdr.PackStart(navBox)

	randBtn, err := gtk.ButtonNewFromIconName("media-playlist-shuffle-symbolic", gtk.ICON_SIZE_LARGE_TOOLBAR)
	if err != nil {
		return nil, err
	}
	randBtn.Connect("clicked", w.RandomComic)
	w.hdr.PackStart(randBtn)

	menu, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	cogImg, err := gtk.ImageNewFromIconName("open-menu", gtk.ICON_SIZE_LARGE_TOOLBAR)
	if err != nil {
		return nil, err
	}
	menu.SetImage(cogImg)

	// Create the cog menu.
	popover, err := gtk.PopoverNew(menu)
	if err != nil {
		return nil, err
	}
	menu.SetPopover(popover)

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		return nil, err
	}
	// TODO: this should be a GtkModelButton, but gotk3 doesn't support
	// it yet.
	menuProp, err := gtk.ButtonNewWithLabel("Properties")
	if err != nil {
		return nil, err
	}
	menuProp.SetRelief(gtk.RELIEF_NONE)
	menuProp.Connect("clicked", w.ShowProperties)
	box.Add(menuProp)
	menuSep, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	box.Add(menuSep)
	// TODO: this should be a GtkModelButton, but gotk3 doesn't support
	// it yet.
	menuAbout, err := gtk.ButtonNewWithLabel("About")
	if err != nil {
		return nil, err
	}
	menuAbout.SetRelief(gtk.RELIEF_NONE)
	menuAbout.Connect("clicked", showAboutDialog)
	box.Add(menuAbout)
	box.ShowAll()
	popover.Add(box)

	w.hdr.PackEnd(menu)

	searchBtn, err := gtk.ButtonNewFromIconName("edit-find", gtk.ICON_SIZE_LARGE_TOOLBAR)
	if err != nil {
		return nil, err
	}
	w.hdr.PackEnd(searchBtn)

	w.hdr.ShowAll()
	w.win.SetTitlebar(w.hdr)

	// Create main part of window.
	scwin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}

	w.img, err = gtk.ImageNewFromIconName("emblem-synchronizing-symbolic", gtk.ICON_SIZE_DIALOG)
	if err != nil {
		return nil, err
	}
	scwin.Add(w.img)
	scwin.ShowAll()
	w.win.Add(scwin)

	return w, nil
}

// PreviousComic sets the current comic to the previous comic.
func (w *Window) PreviousComic() {
	err := w.SetComic(w.comic.Num - 1)
	if err != nil {
		log.Print(err)
	}
}

// NextComic sets the current comic to the next comic.
func (w *Window) NextComic() {
	err := w.SetComic(w.comic.Num + 1)
	if err != nil {
		log.Print(err)
	}
}

// RandomComic sets the current comic to a random comic.
func (w *Window) RandomComic() {
	c, err := getNewestComicInfo()
	if err != nil {
		log.Print(err)
		return
	}
	err = w.SetComic(rand.Intn(c.Num) + 1)
	if err != nil {
		log.Print(err)
	}
}

// SetComic sets the current comic to the given comic.
func (w *Window) SetComic(n int) error {
	var c *xkcd.Comic
	var err error
	if n == 0 {
		c, err = getNewestComicInfo()
		if err != nil {
			return err
		}
	} else {
		c, err = getComicInfo(n)
		if err != nil {
			return err
		}
	}
	w.comic = c

	imgPath, err := getComicImage(w.comic.Num)
	if err != nil {
		log.Printf("error downloading comic: %v", w.comic.Num)
	}
	w.hdr.SetTitle(w.comic.Title)
	w.hdr.SetSubtitle(strconv.Itoa(w.comic.Num))
	w.img.SetFromFile(imgPath)
	w.img.SetTooltipText(w.comic.Alt)

	// Enable/disable previous button.
	if w.comic.Num > 1 {
		w.previous.SetSensitive(true)
	} else {
		w.previous.SetSensitive(false)
	}

	// Enable/disable next button.
	newest, err := getNewestComicInfo()
	if err != nil {
		return err
	}
	if w.comic.Num < newest.Num {
		w.next.SetSensitive(true)
	} else {
		w.next.SetSensitive(false)
	}

	return nil
}

func (w *Window) ShowProperties() {
	pd, err := NewPropertiesDialog(w.win, w.comic)
	if err != nil {
		log.Print(err)
		return
	}
	pd.Present()
}
