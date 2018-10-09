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
	win   *gtk.ApplicationWindow
	state WindowState

	comic      *xkcd.Comic
	comicMutex *sync.Mutex

	actions     map[string]*glib.SimpleAction
	actionGroup *glib.SimpleActionGroup

	hdr *gtk.HeaderBar
	img *gtk.Image

	previous *gtk.Button
	next     *gtk.Button
	rand     *gtk.Button
	search   *gtk.MenuButton
	menu     *gtk.MenuButton

	searchEntry   *gtk.SearchEntry
	searchResults *gtk.Box

	properties *PropertiesDialog
}

// NewWindow creates a new xkcd viewer window.
func NewWindow(app *Application) (*Window, error) {
	var err error

	w := new(Window)

	w.comic = &xkcd.Comic{Title: appName}
	w.comicMutex = new(sync.Mutex)

	w.win, err = gtk.ApplicationWindowNew(app.application)
	if err != nil {
		return nil, err
	}

	actionFuncs := map[string]interface{}{
		"explain":         w.Explain,
		"goto-newest":     w.GotoNewest,
		"open-link":       w.OpenLink,
		"show-properties": w.ShowProperties,
	}

	w.actions = make(map[string]*glib.SimpleAction)
	w.actionGroup = glib.SimpleActionGroupNew()
	for name, function := range actionFuncs {
		action := glib.SimpleActionNew(name, nil)
		action.Connect("activate", function)

		w.actions[name] = action
		w.actionGroup.AddAction(action)
	}
	w.win.InsertActionGroup("win", w.actionGroup)

	// If the gtk theme changes, we might want to adjust our styling.
	w.win.Window.Connect("style-updated", w.StyleUpdated)

	// Create HeaderBar
	w.hdr, err = gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}
	w.hdr.SetTitle(appName)
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
	w.previous.SetTooltipText("Go to the previous comic")
	w.previous.Connect("clicked", w.PreviousComic)
	navBox.Add(w.previous)

	w.next, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	w.next.SetTooltipText("Go to the next comic")
	w.next.Connect("clicked", w.NextComic)
	navBox.Add(w.next)

	w.hdr.PackStart(navBox)

	w.rand, err = gtk.ButtonNewWithLabel("Random")
	if err != nil {
		return nil, err
	}
	w.rand.SetTooltipText("Go to a random comic")
	w.rand.Connect("clicked", w.RandomComic)
	w.hdr.PackStart(w.rand)

	// Create the menu
	w.menu, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	w.menu.SetTooltipText("Menu")

	menuSection1 := glib.MenuNew()
	menuSection1.Append("Open Link", "win.open-link")
	menuSection1.Append("Explain", "win.explain")
	menuSection1.Append("Properties", "win.show-properties")
	menuSection2 := glib.MenuNew()
	menuSection2.Append("Go to Newest Comic", "win.goto-newest")
	menuSection2.Append("New Window", "app.new-window")
	menuSection3 := glib.MenuNew()
	menuSection3.Append("what if?", "app.open-what-if")
	menuSection3.Append("xkcd blog", "app.open-blog")
	menuSection3.Append("xkcd store", "app.open-store")
	menuSection4 := glib.MenuNew()
	menuSection4.Append("About "+appName, "app.show-about")
	menu := glib.MenuNew()
	menu.AppendSectionWithoutLabel(&menuSection1.MenuModel)
	menu.AppendSectionWithoutLabel(&menuSection2.MenuModel)
	menu.AppendSectionWithoutLabel(&menuSection3.MenuModel)
	menu.AppendSectionWithoutLabel(&menuSection4.MenuModel)

	w.menu.SetMenuModel(&menu.MenuModel)
	w.hdr.PackEnd(w.menu)

	// Create the search menu
	w.search, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	w.search.SetTooltipText("Search")
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
	imgScroller, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	imgScroller.SetSizeRequest(400, 300)

	imgContext, err := imgScroller.GetStyleContext()
	if err != nil {
		return nil, err
	}
	imgContext.AddClass("comic-container")

	w.img, err = gtk.ImageNew()
	if err != nil {
		return nil, err
	}
	imgScroller.Add(w.img)
	imgScroller.ShowAll()
	w.win.Add(imgScroller)

	// Recall our window state.
	w.state.ReadFile(filepath.Join(CacheDir(), "state"))
	if w.state.Maximized {
		w.win.Maximize()
	} else {
		w.win.Resize(w.state.Width, w.state.Height)
		if w.state.PositionX != 0 && w.state.PositionY != 0 {
			w.win.Move(w.state.PositionX, w.state.PositionY)
		}
	}
	if w.state.PropertiesVisible {
		if w.properties == nil {
			w.properties, err = NewPropertiesDialog(w)
			if err != nil {
				return nil, err
			}
		}
		w.properties.Present()
	}
	w.SetComic(w.state.ComicNumber)

	// If the gtk window state changes, we want to update our internal
	// window state.
	w.win.Window.Connect("size-allocate", w.StateChanged)
	w.win.Window.Connect("window-state-event", w.StateChanged)

	// If the window is closed, we want to write our state to disk.
	w.win.Window.Connect("delete-event", w.SaveState)

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
	w.state.ComicNumber = n

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
		w.actions["open-link"].SetEnabled(false)
	} else {
		w.actions["open-link"].SetEnabled(true)
	}

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

// GotoNewest checks for a new comic and then shows the newest comic to
// the user.
func (w *Window) GotoNewest() {
	// Make it clear that we are checking for a new comic.
	w.hdr.SetTitle("Checking for new comic...")
	// Close the menu.
	w.menu.GetPopup().Popdown()

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
	err := open.Start(fmt.Sprintf("https://www.explainxkcd.com/%v/", w.comic.Num))
	if err != nil {
		log.Print(err)
	}
}

// OpenLink opens the comic's Link in the user's web browser..
func (w *Window) OpenLink() {
	err := open.Start(w.comic.Link)
	if err != nil {
		log.Print(err)
	}
}
