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

const (
	whatifLink = "https://what-if.xkcd.com/"
	blogLink   = "https://blog.xkcd.com/"
	storeLink  = "https://store.xkcd.com/"
)

// Window is the main application window.
type Window struct {
	state *WindowState

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

	// If the gtk theme changes, we might want to adjust our styling.
	w.win.Connect("style-updated", w.StyleUpdated)

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

	menu, err := gtk.MenuNew()
	if err != nil {
		return nil, err
	}
	menu.SetHAlign(gtk.ALIGN_END)

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
	menuWebsiteWhatIf, err := gtk.MenuItemNewWithLabel("what if?")
	if err != nil {
		return nil, err
	}
	menuWebsiteWhatIf.Connect("activate", OpenURL, whatifLink)
	menuWebsiteWhatIf.SetTooltipText(whatifLink)
	menu.Add(menuWebsiteWhatIf)
	menuWebsiteBlog, err := gtk.MenuItemNewWithLabel("xkcd blog")
	if err != nil {
		return nil, err
	}
	menuWebsiteBlog.Connect("activate", OpenURL, blogLink)
	menuWebsiteBlog.SetTooltipText(blogLink)
	menu.Add(menuWebsiteBlog)
	menuWebsiteStore, err := gtk.MenuItemNewWithLabel("xkcd store")
	if err != nil {
		return nil, err
	}
	menuWebsiteStore.Connect("activate", OpenURL, storeLink)
	menuWebsiteStore.SetTooltipText(storeLink)
	menu.Add(menuWebsiteStore)
	w.menu.SetPopup(menu)
	menu.ShowAll()

	w.hdr.PackEnd(w.menu)

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
	mainScroller, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	mainScroller.SetSizeRequest(400, 300)

	w.img, err = gtk.ImageNew()
	if err != nil {
		return nil, err
	}
	mainScroller.Add(w.img)
	mainScroller.ShowAll()
	w.win.Add(mainScroller)

	// Recall our window state.
	w.state = new(WindowState)
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
	w.win.Connect("size-allocate", w.StateChanged)
	w.win.Connect("window-state-event", w.StateChanged)

	// If the window is closed, we want to write our state to disk.
	w.win.Connect("delete-event", w.SaveState)

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
	w.menuOpenLink.SetTooltipText(w.comic.Link)
	if w.comic.Link == "" {
		w.menuOpenLink.SetSensitive(false)
	} else {
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

var (
	// largeToolbarThemes is the list of gtk themes for which we should use
	// large toolbar buttons.
	largeToolbarThemes = []string{"elementary", "win32"}

	// symbolicIconThemes is the list of gtk themes for which we should use
	// symbolic icons.
	symbolicIconThemes = []string{"Adwaita"}
)

// StyleUpdated is called when the style of our gtk window is updated.
func (w *Window) StyleUpdated() {
	// First, lets find out what GTK theme we are using.
	themeName := os.Getenv("GTK_THEME")
	if themeName == "" {
		// The theme is not being set by the environment, so lets ask
		// GTK what theme it is going to use.
		settings, err := gtk.SettingsGetDefault()
		if err != nil {
			log.Print(err)
		} else {
			// settings.GetProperty returns an interface{}, we will convert
			// it to a string in a moment.
			themeNameIface, err := settings.GetProperty("gtk-theme-name")
			if err != nil {
				log.Print(err)
			} else {
				themeNameStr, ok := themeNameIface.(string)
				if ok {
					themeName = themeNameStr
				}
			}
		}
	}

	// The default size for our headerbar buttons is small.
	headerBarIconSize := gtk.ICON_SIZE_SMALL_TOOLBAR
	for _, largeToolbarTheme := range largeToolbarThemes {
		if themeName == largeToolbarTheme {
			headerBarIconSize = gtk.ICON_SIZE_LARGE_TOOLBAR
		}
	}

	// Should we use symbolic icons?
	useSymbolicIcons := false
	for _, symbolicIconTheme := range symbolicIconThemes {
		if themeName == symbolicIconTheme {
			useSymbolicIcons = true
		}
	}
	// we will call icon() to automatically add -symbolic if needed.
	icon := func(s string) string {
		if useSymbolicIcons {
			return fmt.Sprint(s, "-symbolic")
		}
		return s
	}

	nextImg, err := gtk.ImageNewFromIconName(icon("go-next"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		w.next.SetImage(nextImg)
	}

	previousImg, err := gtk.ImageNewFromIconName(icon("go-previous"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		w.previous.SetImage(previousImg)
	}

	searchImg, err := gtk.ImageNewFromIconName(icon("edit-find"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		w.search.SetImage(searchImg)
	}

	menuImg, err := gtk.ImageNewFromIconName(icon("open-menu"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		w.menu.SetImage(menuImg)
	}
}
