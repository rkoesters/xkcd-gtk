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
	"sync"
)

// Window is the main application window.
type Window struct {
	comic      *xkcd.Comic
	comicMutex *sync.Mutex

	win *gtk.ApplicationWindow
	hdr *gtk.HeaderBar
	img *gtk.Image

	previous *gtk.Button
	next     *gtk.Button
	rand     *gtk.Button
	search   *gtk.MenuButton
	menu     *gtk.MenuButton

	gotoDialog *GotoDialog
	properties *PropertiesDialog

	searchEntry   *gtk.SearchEntry
	searchResults *gtk.Box

	menuExplain  *gtk.MenuItem
	menuOpenLink *gtk.MenuItem
}

// NewWindow creates a new XKCD viewer window.
func NewWindow(app *Application) (*Window, error) {
	var err error

	w := new(Window)

	w.comic = &xkcd.Comic{Title: "XKCD Viewer"}
	w.comicMutex = new(sync.Mutex)

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

	navBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}
	navBoxStyleContext, err := navBox.GetStyleContext()
	if err != nil {
		return nil, err
	}
	navBoxStyleContext.AddClass("linked")

	w.previous, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	w.previous.Connect("clicked", w.PreviousComic)
	navBox.Add(w.previous)

	w.next, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	w.next.Connect("clicked", w.NextComic)
	navBox.Add(w.next)

	w.hdr.PackStart(navBox)

	w.rand, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	w.rand.Connect("clicked", w.RandomComic)
	w.hdr.PackStart(w.rand)

	// Create the menu
	w.menu, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}

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
	w.menu.SetPopup(menu)
	menu.ShowAll()

	w.hdr.PackEnd(w.menu)

	w.search, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	w.hdr.PackEnd(w.search)

	searchPopover, err := gtk.PopoverNew(w.search)
	if err != nil {
		return nil, err
	}
	w.search.SetPopover(searchPopover)

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
	w.searchEntry.Connect("search-changed", w.Search)
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

	// If the gtk theme changes, we might want to adjust our styling.
	w.win.Connect("style-updated", w.StyleUpdatedEvent)
	w.StyleUpdatedEvent()

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

	go func() {
		var err error

		// Make sure we are the only ones changing w.comic.
		w.comicMutex.Lock()
		defer w.comicMutex.Unlock()

		w.comic, err = GetComicInfo(n)
		if err != nil {
			log.Printf("error downloading comic info: %v", n)
		} else {
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

	// If the comic has a link, lets give the option of visiting it.
	if w.comic.Link == "" {
		w.menuOpenLink.SetTooltipText("")
		w.menuOpenLink.SetSensitive(false)
	} else {
		w.menuOpenLink.SetTooltipText(w.comic.Link)
		w.menuOpenLink.SetSensitive(true)
	}
	w.menuExplain.SetTooltipText(explainURL(w.comic.Num))

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

// ShowProperties presents the properties dialog to the user. If the
// dialog doesn't exist yet, we create it.
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

// ShowGoto presents the goto dialog to the user. If the dialog doesn't
// exist yet, we create it.
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

// GotoNewest checks for a new comic and then shows the newest comic to
// the user.
func (w *Window) GotoNewest() {
	// Make it clear that we are checking for a new comic.
	w.hdr.SetTitle("Checking for new comic...")

	// Force GetNewestComicInfo to check for a new comic.
	cachedNewestComic = nil
	newestComic, err := GetNewestComicInfo()
	if err != nil {
		log.Print(err)
	}
	w.SetComic(newestComic.Num)
}

// Explain opens a link to explainxkcd.com in the user's web browser.
func (w *Window) Explain() {
	err := open.Start(explainURL(w.comic.Num))
	if err != nil {
		log.Print(err)
	}
}

func explainURL(n int) string {
	return fmt.Sprintf("https://www.explainxkcd.com/%v/", n)
}

// OpenLink opens the comic's Link in the user's web browser..
func (w *Window) OpenLink() {
	err := open.Start(w.comic.Link)
	if err != nil {
		log.Print(err)
	}
}

// DeleteEvent gets called when our window gets deleted, and we want to
// save our window state for next time.
func (w *Window) DeleteEvent() {
	// Remember what comic we were viewing.
	ws := NewWindowState(w)
	err := ws.WriteFile(filepath.Join(CacheDir(), "state"))
	if err != nil {
		log.Print(err)
	}
}

// StyleUpdatedEvent is called when the style of our gtk window is
// updated.
func (w *Window) StyleUpdatedEvent() {
	log.Print("StyleUpdateEvent()")

	headerBarIconSize := gtk.ICON_SIZE_SMALL_TOOLBAR

	// First, lets find out what theme we are using and use that to
	// decide the size of HeaderBar icons.
	settings, err := gtk.SettingsGetDefault()
	if err != nil {
		log.Print(err)
	} else {
		themeName, err := settings.GetProperty("gtk-theme-name")
		if err != nil {
			log.Print(err)
		}
		themeNameString, ok := themeName.(string)
		if ok && err == nil {
			if themeNameString == "elementary" {
				headerBarIconSize = gtk.ICON_SIZE_LARGE_TOOLBAR
			}
		}
	}

	nextImg, err := gtk.ImageNewFromIconName("go-next-symbolic", headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		w.next.SetImage(nextImg)
	}

	previousImg, err := gtk.ImageNewFromIconName("go-previous-symbolic", headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		w.previous.SetImage(previousImg)
	}

	randImg, err := gtk.ImageNewFromIconName("media-playlist-shuffle-symbolic", headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		w.rand.SetImage(randImg)
	}

	searchImg, err := gtk.ImageNewFromIconName("edit-find", headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		w.search.SetImage(searchImg)
	}

	menuImg, err := gtk.ImageNewFromIconName("open-menu", headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		w.menu.SetImage(menuImg)
	}
}
