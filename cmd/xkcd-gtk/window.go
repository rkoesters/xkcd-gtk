package main

import (
	"fmt"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/style"
	"github.com/rkoesters/xkcd-gtk/internal/widget"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Window is the main application window.
type Window struct {
	app    *Application
	window *gtk.ApplicationWindow
	state  widget.WindowState

	comic      *xkcd.Comic
	comicMutex sync.RWMutex

	actions map[string]*glib.SimpleAction
	accels  *gtk.AccelGroup

	header *gtk.HeaderBar

	navigationBar *widget.NavigationBar
	searchMenu    *widget.SearchMenu
	bookmarksMenu *widget.BookmarksMenu
	windowMenu    *widget.WindowMenu

	comicContainer *widget.ImageViewer

	properties *PropertiesDialog
}

// NewWindow creates a new xkcd viewer window.
func NewWindow(app *Application) (*Window, error) {
	var err error

	win := &Window{
		app: app,
	}

	win.window, err = gtk.ApplicationWindowNew(app.application)
	if err != nil {
		return nil, err
	}

	win.comicMutex.Lock()
	win.comic = &xkcd.Comic{Title: appName}
	win.comicMutex.Unlock()

	// Initialize our window actions.
	win.actions = make(map[string]*glib.SimpleAction)
	registerAction := func(name string, fn interface{}) {
		action := glib.SimpleActionNew(name, nil)
		action.Connect("activate", fn)

		win.actions[name] = action
		win.window.AddAction(action)
	}

	registerAction("bookmark-new", func() { win.bookmarksMenu.AddBookmark() })
	registerAction("bookmark-remove", func() { win.bookmarksMenu.RemoveBookmark() })
	registerAction("explain", win.Explain)
	registerAction("first-comic", win.FirstComic)
	registerAction("newest-comic", win.NewestComic)
	registerAction("next-comic", win.NextComic)
	registerAction("open-link", win.OpenLink)
	registerAction("previous-comic", win.PreviousComic)
	registerAction("random-comic", win.RandomComic)
	registerAction("show-properties", win.ShowProperties)

	// Initialize our window accelerators.
	win.accels, err = gtk.AccelGroupNew()
	if err != nil {
		return nil, err
	}
	win.window.AddAccelGroup(win.accels)

	// If the gtk theme changes, we might want to adjust our styling.
	win.window.Connect("style-updated", win.StyleUpdated)

	darkModeSignal := app.gtkSettings.Connect("notify::gtk-application-prefer-dark-theme", win.DrawComic)
	win.window.Connect("delete-event", func() {
		app.gtkSettings.HandlerDisconnect(darkModeSignal)
	})

	// If the window is closed, we want to write our state to disk.
	win.window.Connect("delete-event", func() {
		if win.properties == nil {
			win.state.SaveState(win.window, nil)
		} else {
			win.state.SaveState(win.window, win.properties.dialog)
		}
	})

	// When gtk destroys the window, we want to clean up.
	win.window.Connect("destroy", win.Destroy)

	// Create HeaderBar
	win.header, err = gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}
	win.header.SetTitle(appName)
	win.header.SetShowCloseButton(true)

	// Create navigation buttons
	win.navigationBar, err = widget.NewNavigationBar(win.actions, win.accels)
	if err != nil {
		return nil, err
	}
	win.header.PackStart(win.navigationBar.IWidget())

	// Create the window menu.
	win.windowMenu, err = widget.NewWindowMenu(app.application.PrefersAppMenu(), win.actions, win.accels, win.ShowProperties)
	if err != nil {
		return nil, err
	}
	win.header.PackEnd(win.windowMenu.IWidget())

	// Create the bookmarks menu.
	win.bookmarksMenu, err = widget.NewBookmarksMenu(&win.app.bookmarks, win.window, &win.state, win.actions, win.accels, win.SetComic)
	if err != nil {
		return nil, err
	}
	win.header.PackEnd(win.bookmarksMenu.IWidget())

	// Create the search menu.
	win.searchMenu, err = widget.NewSearchMenu(win.actions, win.accels, win.SetComic)
	if err != nil {
		return nil, err
	}
	win.header.PackEnd(win.searchMenu.IWidget())

	win.header.ShowAll()
	win.window.SetTitlebar(win.header)

	// Create main part of window.
	win.comicContainer, err = widget.NewImageViewer()
	if err != nil {
		return nil, err
	}
	win.comicContainer.Show()
	win.window.Add(win.comicContainer.IWidget())

	// Recall our window state.
	win.state.LoadState()
	win.window.Resize(win.state.Width, win.state.Height)
	if win.state.PositionX != 0 && win.state.PositionY != 0 {
		win.window.Move(win.state.PositionX, win.state.PositionY)
	}
	if win.state.Maximized {
		win.window.Maximize()
	}
	if win.state.PropertiesVisible {
		win.ShowProperties()
	}
	win.SetComic(win.state.ComicNumber)

	return win, nil
}

// StyleUpdated is called when the style of our gtk window is updated.
func (win *Window) StyleUpdated() {
	// Reload app CSS, if needed.
	darkMode := win.app.DarkMode()
	err := style.UpdateCSS(darkMode)
	if err != nil {
		log.Printf("style.UpdateCSS(darkMode=%v) -> %v", darkMode, err)
	}

	// What GTK theme we are using?
	themeName := os.Getenv("GTK_THEME")
	if themeName == "" {
		// The theme is not being set by the environment, so lets ask
		// GTK what theme it is going to use.
		themeNameIface, err := win.app.gtkSettings.GetProperty("gtk-theme-name")
		if err != nil {
			log.Print(err)
		} else {
			themeNameStr, ok := themeNameIface.(string)
			if ok {
				themeName = themeNameStr
			}
		}
	}

	// The default size for our headerbar buttons is small.
	headerBarIconSize := gtk.ICON_SIZE_SMALL_TOOLBAR
	if style.IsLargeToolbarTheme(themeName) {
		headerBarIconSize = gtk.ICON_SIZE_LARGE_TOOLBAR
	}

	useSymbolicIcons := style.IsSymbolicIconTheme(themeName)

	// We will call icon() to automatically add -symbolic if needed.
	icon := func(s string) string {
		if useSymbolicIcons {
			return s + "-symbolic"
		}
		return s
	}

	firstImg, err := gtk.ImageNewFromIconName("go-first-symbolic", headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.navigationBar.SetFirstButtonImage(firstImg)
	}

	previousImg, err := gtk.ImageNewFromIconName("go-previous-symbolic", headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.navigationBar.SetPreviousButtonImage(previousImg)
	}

	randomImg, err := gtk.ImageNewFromIconName("media-playlist-shuffle-symbolic", headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.navigationBar.SetRandomButtonImage(randomImg)
	}

	nextImg, err := gtk.ImageNewFromIconName("go-next-symbolic", headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.navigationBar.SetNextButtonImage(nextImg)
	}

	newestImg, err := gtk.ImageNewFromIconName("go-last-symbolic", headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.navigationBar.SetNewestButtonImage(newestImg)
	}

	bookmarksImg, err := gtk.ImageNewFromIconName(icon("user-bookmarks"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.bookmarksMenu.IWidget().(*gtk.MenuButton).SetImage(bookmarksImg)
	}

	searchImg, err := gtk.ImageNewFromIconName(icon("edit-find"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.searchMenu.IWidget().(*gtk.MenuButton).SetImage(searchImg)
	}

	menuImg, err := gtk.ImageNewFromIconName(icon("open-menu"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.windowMenu.IWidget().(*gtk.MenuButton).SetImage(menuImg)
	}
}

// FirstComic goes to the first comic.
func (win *Window) FirstComic() {
	win.SetComic(1)
}

// PreviousComic sets the current comic to the previous comic.
func (win *Window) PreviousComic() {
	win.SetComic(win.comicNumber() - 1)
}

// NextComic sets the current comic to the next comic.
func (win *Window) NextComic() {
	win.SetComic(win.comicNumber() + 1)
}

// NewestComic checks for a new comic and then shows the newest comic to the
// user.
func (win *Window) NewestComic() {
	// Make it clear that we are checking for a new comic.
	win.header.SetTitle(l("Checking for new comic..."))

	win.ShowLoading()
	newestComic, err := cache.CheckForNewestComicInfo()
	if err != nil {
		log.Print("error jumping to newest comic: ", err)
	}

	win.SetComic(newestComic.Num)
}

// RandomComic sets the current comic to a random comic.
func (win *Window) RandomComic() {
	today := time.Now()
	if today.Month() == time.April && today.Day() == 1 {
		win.SetComic(4) // chosen by fair dice roll.
		return          // guaranteed to be random.
	}

	win.ShowLoading()
	newestComic, _ := cache.NewestComicInfoFromCache()
	if newestComic.Num <= 0 {
		win.SetComic(newestComic.Num)
	} else {
		win.SetComic(rand.Intn(newestComic.Num) + 1)
	}
}

// SetComic sets the current comic to the given comic.
func (win *Window) SetComic(n int) {
	win.state.ComicNumber = n

	// Make it clear that we are loading a comic.
	win.ShowLoading()
	win.header.SetSubtitle(strconv.Itoa(n))

	// Update UI to reflect new current comic.
	win.updateNextPreviousButtonStatus()
	win.bookmarksMenu.UpdateBookmarkButton()

	go func() {
		var err error

		// Make sure we are the only ones changing win.comic.
		win.comicMutex.Lock()
		defer win.comicMutex.Unlock()

		win.comic, err = cache.ComicInfo(n)
		if err != nil {
			log.Printf("error downloading comic info: %v", n)
		} else {
			_, err = os.Stat(cache.ComicImagePath(n))
			if os.IsNotExist(err) {
				err = cache.DownloadComicImage(n)
				if err != nil {
					// We can be sneaky, we use SafeTitle
					// for window title, but we can leave
					// Title alone so the properties dialog
					// can still be correct.
					win.comic.SafeTitle = l("Connect to the internet to download comic image")
				}
			} else if err != nil {
				log.Print("error finding comic image in cache: ", err)
			}
		}

		// Add the DisplayComic function to the event loop so our UI
		// gets updated with the new comic.
		glib.IdleAdd(win.DisplayComic)
	}()
}

// ShowLoading makes the window indicate that it is loading.
func (win *Window) ShowLoading() {
	win.header.SetTitle(l("Loading comic..."))
	win.comicContainer.ShowLoadingScreen()
}

// DisplayComic updates the UI to show the contents of win.comic.
func (win *Window) DisplayComic() {
	win.comicMutex.RLock()
	defer win.comicMutex.RUnlock()

	win.header.SetTitle(win.comic.SafeTitle)
	win.header.SetSubtitle(strconv.Itoa(win.comic.Num))
	win.comicContainer.SetTooltipText(win.comic.Alt)
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

// DrawComic draws the comic and inverts it if we are in dark mode.
func (win *Window) DrawComic() {
	win.comicContainer.SetFromFile(cache.ComicImagePath(win.comicNumber()), win.app.DarkMode())
}

func (win *Window) updateNextPreviousButtonStatus() {
	n := win.comicNumber()

	// Enable/disable previous button.
	if n > 1 {
		win.actions["previous-comic"].SetEnabled(true)
	} else {
		win.actions["previous-comic"].SetEnabled(false)
	}

	// Enable/disable next button with data from cache.
	newest, _ := cache.NewestComicInfoFromCache()
	if n < newest.Num {
		win.actions["next-comic"].SetEnabled(true)
	} else {
		win.actions["next-comic"].SetEnabled(false)
	}

	// Asynchronously enable/disable next button with data from internet.
	go func() {
		newest, _ := cache.CheckForNewestComicInfo()
		if win.comicNumber() < newest.Num {
			glib.IdleAdd(func() {
				win.actions["next-comic"].SetEnabled(true)
			})
		} else {
			glib.IdleAdd(func() {
				win.actions["next-comic"].SetEnabled(false)
			})
		}
	}()
}

// Explain opens a link to explainxkcd.com in the user's web browser.
func (win *Window) Explain() {
	openURL(fmt.Sprintf("https://www.explainxkcd.com/%v/", win.comicNumber()))
}

// OpenLink opens the comic's Link in the user's web browser.
func (win *Window) OpenLink() {
	win.comicMutex.RLock()
	link := win.comic.Link
	win.comicMutex.RUnlock()

	openURL(link)
}

// comicNumber returns the number of the current comic in a thread-safe way. Do
// not call this method if you already hold win.comicMutex.
func (win *Window) comicNumber() int {
	win.comicMutex.RLock()
	defer win.comicMutex.RUnlock()

	return win.comic.Num
}

// Destroy releases all references in the Window struct. This is needed to
// mitigate a memory leak when closing windows.
func (win *Window) Destroy() {
	win.app = nil
	win.window = nil

	win.comic = nil

	win.actions = nil
	win.accels = nil

	win.header = nil

	win.navigationBar.Destroy()
	win.navigationBar = nil

	win.searchMenu.Destroy()
	win.searchMenu = nil

	win.bookmarksMenu.Destroy()
	win.bookmarksMenu = nil

	win.windowMenu.Destroy()
	win.windowMenu = nil

	win.comicContainer.Destroy()
	win.comicContainer = nil

	win.properties = nil

	runtime.GC()
}
