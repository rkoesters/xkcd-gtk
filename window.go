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

	w.previous, err = gtk.ButtonNewFromIconName("go-previous-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
	if err != nil {
		return nil, err
	}
	w.previous.Connect("clicked", w.PreviousComic)
	navBox.Add(w.previous)

	w.next, err = gtk.ButtonNewFromIconName("go-next-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
	if err != nil {
		return nil, err
	}
	w.next.Connect("clicked", w.NextComic)
	navBox.Add(w.next)

	w.hdr.PackStart(navBox)

	randBtn, err := gtk.ButtonNewFromIconName("media-playlist-shuffle-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
	if err != nil {
		return nil, err
	}
	randBtn.Connect("clicked", w.RandomComic)
	w.hdr.PackStart(randBtn)

	menuBtn, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	cogImg, err := gtk.ImageNewFromIconName("open-menu", gtk.ICON_SIZE_SMALL_TOOLBAR)
	if err != nil {
		return nil, err
	}
	menuBtn.SetImage(cogImg)

	menu, err := gtk.MenuNew()
	if err != nil {
		return nil, err
	}

	menuGoto, err := gtk.MenuItemNewWithLabel("Go to...")
	if err != nil {
		return nil, err
	}
	menuGoto.Connect("activate", w.ShowGoto)
	menu.Add(menuGoto)
	menuProp, err := gtk.MenuItemNewWithLabel("Properties")
	if err != nil {
		return nil, err
	}
	menuProp.Connect("activate", w.ShowProperties)
	menu.Add(menuProp)
	menuSep, err := gtk.SeparatorMenuItemNew()
	if err != nil {
		return nil, err
	}
	menu.Add(menuSep)
	menuAbout, err := gtk.MenuItemNewWithLabel("About")
	if err != nil {
		return nil, err
	}
	menuAbout.Connect("activate", app.ShowAboutDialog)
	menu.Add(menuAbout)
	menuBtn.SetPopup(menu)
	menu.ShowAll()

	w.hdr.PackEnd(menuBtn)

	searchBtn, err := gtk.ButtonNewFromIconName("edit-find", gtk.ICON_SIZE_SMALL_TOOLBAR)
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
	scwin.SetSizeRequest(400, 300)

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

func (w *Window) ShowGoto() {
	gt, err := NewGoto(w)
	if err != nil {
		log.Print(err)
		return
	}
	defer gt.dialog.Close()
	status := gt.dialog.Run()
	if status == 1 {
		input, err := gt.entry.GetText()
		if err != nil {
			log.Print(err)
			return
		}
		number, err := strconv.Atoi(input)
		if err != nil {
			log.Print(err)
			return
		}
		w.SetComic(number)
	}
}
