package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xdg"
	"github.com/rkoesters/xkcd"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Window is the main application window.
type Window struct {
	app    *Application
	window *gtk.ApplicationWindow
	state  WindowState

	comic      *xkcd.Comic
	comicMutex sync.Mutex

	actions map[string]*glib.SimpleAction
	accels  *gtk.AccelGroup

	header         *gtk.HeaderBar
	comicContainer *gtk.ScrolledWindow
	image          *gtk.Image

	first    *gtk.Button
	previous *gtk.Button
	next     *gtk.Button
	newest   *gtk.Button

	random        *gtk.Button
	search        *gtk.MenuButton
	searchEntry   *gtk.SearchEntry
	searchResults *gtk.Box

	menu *gtk.MenuButton

	properties *PropertiesDialog
}

// NewWindow creates a new xkcd viewer window.
func NewWindow(app *Application) (*Window, error) {
	var err error

	win := new(Window)

	win.app = app

	win.window, err = gtk.ApplicationWindowNew(app.application)
	if err != nil {
		return nil, err
	}

	win.comic = &xkcd.Comic{Title: appName}

	// Initialize our window actions.
	actionFuncs := map[string]interface{}{
		"first-comic":     win.FirstComic,
		"previous-comic":  win.PreviousComic,
		"next-comic":      win.NextComic,
		"newest-comic":    win.NewestComic,
		"random-comic":    win.RandomComic,
		"open-link":       win.OpenLink,
		"explain":         win.Explain,
		"show-properties": win.ShowProperties,
	}

	win.actions = make(map[string]*glib.SimpleAction)
	for name, function := range actionFuncs {
		action := glib.SimpleActionNew(name, nil)
		action.Connect("activate", function)

		win.actions[name] = action
		win.window.AddAction(action)
	}

	// Initialize our window accelerators.
	win.accels, err = gtk.AccelGroupNew()
	if err != nil {
		return nil, err
	}
	win.window.AddAccelGroup(win.accels)

	// If the gtk theme changes, we might want to adjust our styling.
	win.window.Window.Connect("style-updated", win.StyleUpdated)

	gtkSettings, err := gtk.SettingsGetDefault()
	if err != nil {
		return nil, err
	}
	darkModeSignal, err := gtkSettings.Connect("notify::gtk-application-prefer-dark-theme", win.DrawComic)
	if err != nil {
		return nil, err
	}
	win.window.Connect("delete-event", func() {
		gtkSettings.HandlerDisconnect(darkModeSignal)
	})

	// If the gtk window state changes, we want to update our internal
	// window state.
	win.window.Window.Connect("size-allocate", win.StateChanged)
	win.window.Window.Connect("window-state-event", win.StateChanged)

	// If the window is closed, we want to write our state to disk.
	win.window.Window.Connect("delete-event", win.SaveState)

	// Create HeaderBar
	win.header, err = gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}
	win.header.SetTitle(appName)
	win.header.SetShowCloseButton(true)

	// Create navigation buttons
	navBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}
	navBoxStyleContext, err := navBox.GetStyleContext()
	if err != nil {
		return nil, err
	}
	navBoxStyleContext.AddClass("linked")

	win.first, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	win.first.SetTooltipText("Go to the first comic")
	win.first.SetProperty("action-name", "win.first-comic")
	win.first.AddAccelerator("activate", win.accels, gdk.KEY_Home, gdk.GDK_CONTROL_MASK, gtk.ACCEL_VISIBLE)
	navBox.Add(win.first)

	win.previous, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	win.previous.SetTooltipText("Go to the previous comic")
	win.previous.SetProperty("action-name", "win.previous-comic")
	win.previous.AddAccelerator("activate", win.accels, gdk.KEY_Left, gdk.GDK_CONTROL_MASK, gtk.ACCEL_VISIBLE)
	navBox.Add(win.previous)

	win.next, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	win.next.SetTooltipText("Go to the next comic")
	win.next.SetProperty("action-name", "win.next-comic")
	win.next.AddAccelerator("activate", win.accels, gdk.KEY_Right, gdk.GDK_CONTROL_MASK, gtk.ACCEL_VISIBLE)
	navBox.Add(win.next)

	win.newest, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	win.newest.SetTooltipText("Go to the newest comic")
	win.newest.SetProperty("action-name", "win.newest-comic")
	win.newest.AddAccelerator("activate", win.accels, gdk.KEY_End, gdk.GDK_CONTROL_MASK, gtk.ACCEL_VISIBLE)
	navBox.Add(win.newest)

	win.header.PackStart(navBox)

	// Create the menu
	win.menu, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	win.menu.SetTooltipText("Menu")

	menu := glib.MenuNew()

	menuSection1 := glib.MenuNew()
	menuSection1.Append("Open Link", "win.open-link")
	menuSection1.Append("Explain", "win.explain")
	menuSection1.Append("Properties", "win.show-properties")
	menu.AppendSectionWithoutLabel(&menuSection1.MenuModel)
	win.accels.Connect(gdk.KEY_p, gdk.GDK_CONTROL_MASK, gtk.ACCEL_VISIBLE, win.ShowProperties)

	if !app.application.PrefersAppMenu() {
		menuSection2 := glib.MenuNew()
		menuSection2.Append("New Window", "app.new-window")
		menu.AppendSectionWithoutLabel(&menuSection2.MenuModel)

		menuSection3 := glib.MenuNew()
		menuSection3.Append("Toggle Dark Mode", "app.toggle-dark-mode")
		menu.AppendSectionWithoutLabel(&menuSection3.MenuModel)

		menuSection4 := glib.MenuNew()
		menuSection4.Append("what if?", "app.open-what-if")
		menuSection4.Append("xkcd blog", "app.open-blog")
		menuSection4.Append("xkcd store", "app.open-store")
		menu.AppendSectionWithoutLabel(&menuSection4.MenuModel)

		menuSection5 := glib.MenuNew()
		menuSection5.Append("Keyboard Shortcuts", "app.show-shortcuts")
		menuSection5.Append("About "+appName, "app.show-about")
		menu.AppendSectionWithoutLabel(&menuSection5.MenuModel)
	}

	win.menu.SetMenuModel(&menu.MenuModel)
	win.header.PackEnd(win.menu)

	// Create the search menu
	win.search, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	win.search.SetTooltipText("Search")
	win.search.AddAccelerator("activate", win.accels, gdk.KEY_f, gdk.GDK_CONTROL_MASK, gtk.ACCEL_VISIBLE)
	win.header.PackEnd(win.search)

	searchPopover, err := gtk.PopoverNew(win.search)
	if err != nil {
		return nil, err
	}
	win.search.SetPopover(searchPopover)

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		return nil, err
	}
	box.SetMarginTop(12)
	box.SetMarginBottom(12)
	box.SetMarginStart(12)
	box.SetMarginEnd(12)
	win.searchEntry, err = gtk.SearchEntryNew()
	if err != nil {
		return nil, err
	}
	win.searchEntry.Connect("search-changed", win.Search)
	box.Add(win.searchEntry)
	scwin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	box.Add(scwin)
	win.searchResults, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	scwin.Add(win.searchResults)
	scwin.SetSizeRequest(375, 250)
	win.loadSearchResults(nil)
	box.ShowAll()
	searchPopover.Add(box)

	// Create the random button
	win.random, err = gtk.ButtonNewWithLabel("Random")
	if err != nil {
		return nil, err
	}
	win.random.SetTooltipText("Go to a random comic")
	win.random.SetProperty("action-name", "win.random-comic")
	win.random.AddAccelerator("activate", win.accels, gdk.KEY_r, gdk.GDK_CONTROL_MASK, gtk.ACCEL_VISIBLE)
	win.header.PackEnd(win.random)

	win.header.ShowAll()
	win.window.SetTitlebar(win.header)

	// Create main part of window.
	win.comicContainer, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	win.comicContainer.SetSizeRequest(400, 300)

	imageContext, err := win.comicContainer.GetStyleContext()
	if err != nil {
		return nil, err
	}
	imageContext.AddClass("comic-container")

	win.image, err = gtk.ImageNew()
	if err != nil {
		return nil, err
	}
	win.image.SetHAlign(gtk.ALIGN_CENTER)
	win.image.SetVAlign(gtk.ALIGN_CENTER)

	win.comicContainer.Add(win.image)
	win.comicContainer.ShowAll()
	win.window.Add(win.comicContainer)

	// Recall our window state.
	win.state.ReadFile(filepath.Join(CacheDir(), "state"))
	if win.state.Maximized {
		win.window.Maximize()
	} else {
		win.window.Resize(win.state.Width, win.state.Height)
		if win.state.PositionX != 0 && win.state.PositionY != 0 {
			win.window.Move(win.state.PositionX, win.state.PositionY)
		}
	}
	if win.state.PropertiesVisible {
		win.ShowProperties()
	}
	win.SetComic(win.state.ComicNumber)

	return win, nil
}

// FirstComic goes to the first comic.
func (win *Window) FirstComic() {
	win.SetComic(1)
}

// PreviousComic sets the current comic to the previous comic.
func (win *Window) PreviousComic() {
	win.SetComic(win.comic.Num - 1)
}

// NextComic sets the current comic to the next comic.
func (win *Window) NextComic() {
	win.SetComic(win.comic.Num + 1)
}

// NewestComic checks for a new comic and then shows the newest comic to
// the user.
func (win *Window) NewestComic() {
	// Make it clear that we are checking for a new comic.
	win.header.SetTitle("Checking for new comic...")

	// Force GetNewestComicInfo to check for a new comic.
	setCachedNewestComic <- nil
	newestComic, err := GetNewestComicInfo()
	if err != nil {
		log.Print(err)
	}

	win.SetComic(newestComic.Num)
}

// RandomComic sets the current comic to a random comic.
func (win *Window) RandomComic() {
	newestComic, _ := GetNewestComicInfo()
	if newestComic.Num <= 0 {
		win.SetComic(newestComic.Num)
	} else {
		win.SetComic(rand.Intn(newestComic.Num) + 1)
	}
}

// SetComic sets the current comic to the given comic.
func (win *Window) SetComic(n int) {
	// Make it clear that we are loading a comic.
	win.header.SetTitle("Loading comic...")
	win.header.SetSubtitle(strconv.Itoa(n))
	win.updateNextPreviousButtonStatus()
	win.state.ComicNumber = n

	go func() {
		var err error

		// Make sure we are the only ones changing win.comic.
		win.comicMutex.Lock()
		defer win.comicMutex.Unlock()

		win.comic, err = GetComicInfo(n)
		if err != nil {
			log.Printf("error downloading comic info: %v", n)
		} else {
			_, err = os.Stat(getComicImagePath(n))
			if os.IsNotExist(err) {
				err = DownloadComicImage(n)
				if err != nil {
					// We can be sneaky, we use SafeTitle for window
					// title, but we can leave Title alone so the
					// properties dialog can still be correct.
					win.comic.SafeTitle = "Connect to the internet to download comic image"
				}
			} else if err != nil {
				log.Print(err)
			}
		}

		// Add the DisplayComic function to the event loop so our UI
		// gets updated with the new comic.
		glib.IdleAdd(win.DisplayComic)
	}()
}

// DisplayComic updates the UI to show the contents of win.comic
func (win *Window) DisplayComic() {
	win.header.SetTitle(win.comic.SafeTitle)
	win.header.SetSubtitle(strconv.Itoa(win.comic.Num))
	win.image.SetTooltipText(win.comic.Alt)
	win.updateNextPreviousButtonStatus()

	// If the comic has a link, lets give the option of visiting it.
	if win.comic.Link == "" {
		win.actions["open-link"].SetEnabled(false)
	} else {
		win.actions["open-link"].SetEnabled(true)
	}

	if win.properties != nil {
		win.properties.Update()
	}

	win.DrawComic()
}

func (win *Window) updateNextPreviousButtonStatus() {
	// Enable/disable previous button.
	if win.comic.Num > 1 {
		win.actions["previous-comic"].SetEnabled(true)
	} else {
		win.actions["previous-comic"].SetEnabled(false)
	}

	// Enable/disable next button.
	newest, _ := GetNewestComicInfoAsync(func(c *xkcd.Comic, _ error) {
		if c != nil {
			if win.comic.Num < c.Num {
				glib.IdleAdd(func() {
					win.actions["next-comic"].SetEnabled(true)
				})
			} else {
				glib.IdleAdd(func() {
					win.actions["next-comic"].SetEnabled(false)
				})
			}
		}
	})
	if win.comic.Num < newest.Num {
		win.actions["next-comic"].SetEnabled(true)
	} else {
		win.actions["next-comic"].SetEnabled(false)
	}
}

// Explain opens a link to explainxkcd.com in the user's web browser.
func (win *Window) Explain() {
	err := xdg.Open(fmt.Sprintf("https://www.explainxkcd.com/%v/", win.comic.Num))
	if err != nil {
		log.Print(err)
	}
}

// OpenLink opens the comic's Link in the user's web browser.
func (win *Window) OpenLink() {
	err := xdg.Open(win.comic.Link)
	if err != nil {
		log.Print(err)
	}
}
