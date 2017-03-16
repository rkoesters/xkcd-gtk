package main

import (
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd"
	"github.com/skratchdot/open-golang/open"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

// Window is the main application window.
type Window struct {
	comic *xkcd.Comic

	win      *gtk.ApplicationWindow
	hdr      *gtk.HeaderBar
	previous *gtk.Button
	next     *gtk.Button
	rand     *gtk.Button
	img      *gtk.Image

	gotoDialog *GotoDialog
	properties *PropertiesDialog

	searchEntry   *gtk.SearchEntry
	searchResults *gtk.Box

	menuExplain  *gtk.MenuItem
	menuOpenLink *gtk.MenuItem
}

// New creates a new XKCD viewer window.
func NewWindow(app *Application) (*Window, error) {
	var err error

	w := new(Window)

	w.comic = &xkcd.Comic{Title: "XKCD Viewer"}

	w.win, err = gtk.ApplicationWindowNew(app.GtkApp)
	if err != nil {
		return nil, err
	}
	w.win.Connect("delete-event", w.DeleteEvent)
	w.win.SetDefaultSize(1000, 800)

	// Create HeaderBar
	w.hdr, err = gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}
	w.hdr.SetTitle("XKCD Viewer")
	w.hdr.SetShowCloseButton(true)

	headerBarIconSize := lookupHeaderBarIconSize()

	navBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}
	navBoxStyleContext, err := navBox.GetStyleContext()
	if err != nil {
		return nil, err
	}
	navBoxStyleContext.AddClass("linked")

	w.previous, err = gtk.ButtonNewFromIconName("go-previous-symbolic", headerBarIconSize)
	if err != nil {
		return nil, err
	}
	w.previous.Connect("clicked", w.PreviousComic)
	navBox.Add(w.previous)

	w.next, err = gtk.ButtonNewFromIconName("go-next-symbolic", headerBarIconSize)
	if err != nil {
		return nil, err
	}
	w.next.Connect("clicked", w.NextComic)
	navBox.Add(w.next)

	w.hdr.PackStart(navBox)

	w.rand, err = gtk.ButtonNewFromIconName("media-playlist-shuffle-symbolic", headerBarIconSize)
	if err != nil {
		return nil, err
	}
	w.rand.Connect("clicked", w.RandomComic)
	w.hdr.PackStart(w.rand)

	// Create the menu
	menuBtn, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	cogImg, err := gtk.ImageNewFromIconName("open-menu", headerBarIconSize)
	if err != nil {
		return nil, err
	}
	menuBtn.SetImage(cogImg)

	menu, err := gtk.MenuNew()
	if err != nil {
		return nil, err
	}

	w.menuOpenLink, err = gtk.MenuItemNewWithLabel("Open Link")
	if err != nil {
		return nil, err
	}
	w.menuOpenLink.Connect("activate", w.OpenLink)
	menu.Add(w.menuOpenLink)
	w.menuExplain, err = gtk.MenuItemNewWithLabel("Explain")
	if err != nil {
		return nil, err
	}
	w.menuExplain.Connect("activate", w.Explain)
	menu.Add(w.menuExplain)
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
	menuGotoNewest, err := gtk.MenuItemNewWithLabel("Go to Newest Comic")
	if err != nil {
		return nil, err
	}
	menuGotoNewest.Connect("activate", w.GotoNewest)
	menu.Add(menuGotoNewest)
	menuGoto, err := gtk.MenuItemNewWithLabel("Go to...")
	if err != nil {
		return nil, err
	}
	menuGoto.Connect("activate", w.ShowGoto)
	menu.Add(menuGoto)
	menuNewWindow, err := gtk.MenuItemNewWithLabel("New Window")
	if err != nil {
		return nil, err
	}
	menuNewWindow.Connect("activate", app.Activate)
	menu.Add(menuNewWindow)
	menuSep, err = gtk.SeparatorMenuItemNew()
	if err != nil {
		return nil, err
	}
	menu.Add(menuSep)
	menuAbout, err := gtk.MenuItemNewWithLabel("About")
	if err != nil {
		return nil, err
	}
	menuAbout.Connect("activate", ShowAboutDialog)
	menu.Add(menuAbout)
	menuBtn.SetPopup(menu)
	menu.ShowAll()

	w.hdr.PackEnd(menuBtn)

	searchBtn, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	searchImg, err := gtk.ImageNewFromIconName("edit-find", headerBarIconSize)
	if err != nil {
		return nil, err
	}
	searchBtn.SetImage(searchImg)
	w.hdr.PackEnd(searchBtn)

	searchPopover, err := gtk.PopoverNew(searchBtn)
	if err != nil {
		return nil, err
	}
	searchBtn.SetPopover(searchPopover)

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		return nil, err
	}
	box.SetMarginTop(12)
	box.SetMarginBottom(12)
	box.SetMarginStart(12)
	box.SetMarginEnd(12)
	w.searchEntry, err = gtk.SearchEntryNew()
	if err != nil {
		return nil, err
	}
	w.searchEntry.Connect("search-changed", w.UpdateSearch)
	box.Add(w.searchEntry)
	scwin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	box.Add(scwin)
	w.searchResults, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	scwin.Add(w.searchResults)
	scwin.SetSizeRequest(375, 250)
	w.loadSearchResults(nil)
	box.ShowAll()
	searchPopover.Add(box)

	w.hdr.ShowAll()
	w.win.SetTitlebar(w.hdr)

	// Create main part of window.
	searchScroller, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	searchScroller.SetSizeRequest(400, 300)

	w.img, err = gtk.ImageNew()
	if err != nil {
		return nil, err
	}
	searchScroller.Add(w.img)
	searchScroller.ShowAll()
	w.win.Add(searchScroller)

	// Recall our window state.
	ws := NewWindowState(w)
	ws.ReadFile(filepath.Join(CacheDir(), "state"))
	w.win.Resize(ws.Width, ws.Height)
	w.win.Move(ws.PositionX, ws.PositionY)
	if ws.PropertiesVisible {
		if w.properties == nil {
			w.properties, err = NewPropertiesDialog(w)
			if err != nil {
				return nil, err
			}
		}
		w.properties.dialog.Resize(ws.PropertiesWidth, ws.PropertiesHeight)
		w.properties.dialog.Move(ws.PropertiesPositionX, ws.PropertiesPositionY)
		w.properties.Present()
	}
	w.SetComic(ws.ComicNumber)

	return w, nil
}

// PreviousComic sets the current comic to the previous comic.
func (w *Window) PreviousComic() {
	w.SetComic(w.comic.Num - 1)
}

// NextComic sets the current comic to the next comic.
func (w *Window) NextComic() {
	w.SetComic(w.comic.Num + 1)
}

// RandomComic sets the current comic to a random comic.
func (w *Window) RandomComic() {
	newestComic, _ := GetNewestComicInfo()
	if newestComic.Num <= 0 {
		w.SetComic(newestComic.Num)
	} else {
		w.SetComic(rand.Intn(newestComic.Num) + 1)
	}
}

// SetComic sets the current comic to the given comic.
func (w *Window) SetComic(n int) {
	// Make it clear that we are loading a comic.
	w.hdr.SetTitle("Loading comic...")
	w.hdr.SetSubtitle(strconv.Itoa(n))
	w.updateNextPreviousButtonStatus()
	w.rand.SetSensitive(false)

	go func() {
		var err error
		w.comic, err = GetComicInfo(n)
		if err != nil {
			log.Printf("error downloading comic info: %v", n)
		}

		_, err = os.Stat(getComicImagePath(n))
		if os.IsNotExist(err) {
			err = DownloadComicImage(n)
			if err != nil {
				// We can be sneaky, we use SafeTitle for window title,
				// but we can leave Title alone so the properties dialog
				// can still be correct.
				w.comic.SafeTitle = "Connect to the internet to download comic image"
			}
		} else if err != nil {
			log.Print(err)
		}

		// Add the DisplayComic function to the event loop so our UI
		// gets updated with the new comic.
		glib.IdleAdd(w.DisplayComic)
	}()
}

// DisplayComic updates the UI to show the contents of w.comic
func (w *Window) DisplayComic() {
	w.hdr.SetTitle(w.comic.SafeTitle)
	w.hdr.SetSubtitle(strconv.Itoa(w.comic.Num))
	w.img.SetFromFile(getComicImagePath(w.comic.Num))
	w.img.SetTooltipText(w.comic.Alt)
	w.updateNextPreviousButtonStatus()
	w.rand.SetSensitive(true)

	// If the comic has a link, lets give the option of visiting it.
	if w.comic.Link == "" {
		w.menuOpenLink.SetTooltipText("")
		w.menuOpenLink.SetSensitive(false)
	} else {
		w.menuOpenLink.SetTooltipText(w.comic.Link)
		w.menuOpenLink.SetSensitive(true)
	}
	w.menuExplain.SetTooltipText(explainUrl(w.comic.Num))

	if w.properties != nil {
		w.properties.Update()
	}
}

func (w *Window) updateNextPreviousButtonStatus() {
	// Enable/disable previous button.
	if w.comic.Num > 1 {
		w.previous.SetSensitive(true)
	} else {
		w.previous.SetSensitive(false)
	}

	// Enable/disable next button.
	newest, _ := GetNewestComicInfo()
	if w.comic.Num < newest.Num {
		w.next.SetSensitive(true)
	} else {
		w.next.SetSensitive(false)
	}
}

func (w *Window) ShowProperties() {
	var err error
	if w.properties == nil {
		w.properties, err = NewPropertiesDialog(w)
		if err != nil {
			log.Print(err)
			return
		}
	}
	w.properties.Present()
}

func (w *Window) ShowGoto() {
	var err error
	if w.gotoDialog == nil {
		w.gotoDialog, err = NewGotoDialog(w)
		if err != nil {
			log.Print(err)
			return
		}
	}
	w.gotoDialog.Present()
}

func (w *Window) GotoNewest() {
	// Make it clear that we are checking for a new comic.
	w.hdr.SetTitle("Checking for new comic...")
	w.previous.SetSensitive(false)
	w.next.SetSensitive(false)
	w.rand.SetSensitive(false)

	newestComic, err := GetNewestComicInfo()
	if err != nil {
		log.Print(err)
	}
	w.SetComic(newestComic.Num)
}

// Open explainxkcd in web browser.
func (w *Window) Explain() {
	err := open.Start(explainUrl(w.comic.Num))
	if err != nil {
		log.Print(err)
	}
}

func explainUrl(n int) string {
	return fmt.Sprintf("https://www.explainxkcd.com/%v/", n)
}

// Open the comic's link in a web browser.
func (w *Window) OpenLink() {
	err := open.Start(w.comic.Link)
	if err != nil {
		log.Print(err)
	}
}

func (w *Window) DeleteEvent() {
	// Remember what comic we were viewing.
	ws := NewWindowState(w)
	err := ws.WriteFile(filepath.Join(CacheDir(), "state"))
	if err != nil {
		log.Print(err)
	}
}
