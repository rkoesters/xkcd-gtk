package widget

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

// ApplicationWindow is the main application window.
type ApplicationWindow struct {
	*gtk.ApplicationWindow

	app   Application
	state WindowState

	comic      *xkcd.Comic
	comicMutex sync.RWMutex

	actions map[string]*glib.SimpleAction

	header        *gtk.HeaderBar
	navigationBar *NavigationBar
	searchMenu    *SearchMenu
	bookmarksMenu *BookmarksMenu
	windowMenu    *WindowMenu

	comicContainer *ImageViewer

	properties *PropertiesDialog // May be nil.
}

var _ Widget = &ApplicationWindow{}

// NewApplicationWindow creates a new xkcd viewer window.
func NewApplicationWindow(app Application) (*ApplicationWindow, error) {
	super, err := gtk.ApplicationWindowNew(app.GtkApplication())
	if err != nil {
		return nil, err
	}
	win := &ApplicationWindow{
		ApplicationWindow: super,

		app:     app,
		comic:   &xkcd.Comic{Title: AppName()},
		actions: make(map[string]*glib.SimpleAction),
	}

	// Put everything where the user left it.
	win.state.LoadState()

	registerAction := func(name string, fn interface{}) {
		action := glib.SimpleActionNew(name, nil)
		action.Connect("activate", fn)

		win.actions[name] = action
		win.AddAction(action)
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
	registerAction("zoom-in", win.ZoomIn)
	registerAction("zoom-out", win.ZoomOut)
	registerAction("zoom-reset", win.ZoomReset)

	// Initialize our window accelerators.
	accels, err := gtk.AccelGroupNew()
	if err != nil {
		return nil, err
	}
	win.AddAccelGroup(accels)

	accels.Connect(gdk.KEY_plus, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE, win.ZoomIn)
	accels.Connect(gdk.KEY_equal, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE, win.ZoomIn) // without holding shift
	accels.Connect(gdk.KEY_minus, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE, win.ZoomOut)
	accels.Connect(gdk.KEY_0, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE, win.ZoomReset)
	accels.Connect(gdk.KEY_p, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE, win.ShowProperties)
	accels.Connect(gdk.KEY_Return, gdk.MOD1_MASK, gtk.ACCEL_VISIBLE, win.ShowProperties)
	accels.Connect(gdk.KEY_w, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE, win.Close)

	// If the gtk theme changes, we might want to adjust our styling.
	win.Connect("style-updated", win.StyleUpdated)

	darkModeSignal := app.ConnectDarkModeChanged(win.DarkModeChanged)
	win.Connect("delete-event", func() {
		gtks, err := gtk.SettingsGetDefault()
		if err != nil {
			log.Print("error calling gtk.SettingsGetDefault(): ", err)
			return
		}
		gtks.HandlerDisconnect(darkModeSignal)
	})

	// If the window is closed, we want to write our state to disk.
	win.Connect("delete-event", func() {
		win.state.SaveState(win, win.properties)
	})

	// When gtk destroys the window, we want to clean up.
	win.Connect("destroy", win.Dispose)

	// Create image viewing frame
	win.comicContainer, err = NewImageViewer(win.IActionGroup, win.state.ImageScale)
	if err != nil {
		return nil, err
	}
	win.Add(win.comicContainer)
	win.Resize(win.state.Width, win.state.Height)
	if win.state.PositionX != 0 && win.state.PositionY != 0 {
		win.Move(win.state.PositionX, win.state.PositionY)
	}
	if win.state.Maximized {
		win.Maximize()
	}
	if win.state.PropertiesVisible {
		win.ShowProperties()
	}

	// Create HeaderBar
	win.header, err = gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}
	win.header.SetTitle(AppName())
	win.header.SetShowCloseButton(true)

	// Create navigation buttons
	win.navigationBar, err = NewNavigationBar(accels, win.actions, win.comicNumber)
	if err != nil {
		return nil, err
	}
	win.header.PackStart(win.navigationBar)

	// Create the window menu.
	win.windowMenu, err = NewWindowMenu(accels, app.PrefersAppMenu(), app.DarkMode, app.SetDarkMode)
	if err != nil {
		return nil, err
	}
	win.updateZoomButtonStatus()
	win.header.PackEnd(win.windowMenu)

	// Create the bookmarks menu.
	win.bookmarksMenu, err = NewBookmarksMenu(win.app.BookmarksRef(), win, &win.state, win.actions, accels, win.SetComic)
	if err != nil {
		return nil, err
	}
	win.header.PackEnd(win.bookmarksMenu)

	// Create the search menu.
	win.searchMenu, err = NewSearchMenu(accels, win.SetComic)
	if err != nil {
		return nil, err
	}
	win.header.PackEnd(win.searchMenu)

	win.header.ShowAll()
	win.SetTitlebar(win.header)

	win.SetComic(win.state.ComicNumber)

	return win, nil
}

func (win *ApplicationWindow) IsVisible() bool {
	if win == nil {
		return false
	}
	return win.ApplicationWindow.IsVisible()
}

func (win *ApplicationWindow) DarkModeChanged() {
	darkMode := win.app.DarkMode()
	log.Debugf("DarkModeChanged() -> %v", darkMode)
	comicId := win.comicNumber()
	err := win.comicContainer.DrawComic(comicId, darkMode)
	if err != nil {
		log.Print("error calling ImageViewer.DrawComic(id=%v, darkMode=%v) -> %v ", comicId, darkMode, err)
	}
	win.StyleUpdated()
	win.windowMenu.darkModeSwitch.SyncDarkMode(darkMode)
}

// StyleUpdated is called when the style of our gtk window is updated.
func (win *ApplicationWindow) StyleUpdated() {
	themeName, err := win.app.GtkTheme()
	if err != nil {
		log.Print("error querying GTK theme: ", err)
	}

	// The default size for our headerbar buttons is small.
	headerBarIconSize := gtk.ICON_SIZE_SMALL_TOOLBAR
	if style.IsLargeToolbarTheme(themeName) {
		headerBarIconSize = gtk.ICON_SIZE_LARGE_TOOLBAR
	}

	useSymbolicIcons := style.IsSymbolicIconTheme(themeName, win.app.DarkMode())

	// We will call icon() to automatically add -symbolic if needed.
	icon := func(s string) string {
		if useSymbolicIcons {
			return s + "-symbolic"
		}
		return s
	}

	setButtonImageFromIconName := func(icon string, imageSetter func(gtk.IWidget)) {
		img, err := gtk.ImageNewFromIconName(icon, headerBarIconSize)
		if err != nil {
			log.Print(err)
			return
		}
		imageSetter(img)
	}

	setButtonImageFromIconName("go-first-symbolic", win.navigationBar.SetFirstButtonImage)
	setButtonImageFromIconName("go-previous-symbolic", win.navigationBar.SetPreviousButtonImage)
	setButtonImageFromIconName("media-playlist-shuffle-symbolic", win.navigationBar.SetRandomButtonImage)
	setButtonImageFromIconName("go-next-symbolic", win.navigationBar.SetNextButtonImage)
	setButtonImageFromIconName("go-last-symbolic", win.navigationBar.SetNewestButtonImage)
	setButtonImageFromIconName(icon("edit-find"), win.searchMenu.SetImage)
	setButtonImageFromIconName(icon("user-bookmarks"), win.bookmarksMenu.SetImage)
	setButtonImageFromIconName(icon("open-menu"), win.windowMenu.SetImage)

	linked := style.IsLinkedNavButtonsTheme(themeName)
	if err := win.navigationBar.SetLinkedButtons(linked); err != nil {
		log.Print(err)
	}

	compact := style.IsCompactMenuTheme(themeName)
	win.windowMenu.SetCompact(compact)
	win.comicContainer.contextMenu.SetCompact(compact)
}

// FirstComic goes to the first comic.
func (win *ApplicationWindow) FirstComic() {
	win.SetComic(1)
}

// PreviousComic sets the current comic to the previous comic.
func (win *ApplicationWindow) PreviousComic() {
	win.SetComic(win.comicNumber() - 1)
}

// NextComic sets the current comic to the next comic.
func (win *ApplicationWindow) NextComic() {
	win.SetComic(win.comicNumber() + 1)
}

// NewestComic checks for a new comic and then shows the newest comic to the
// user.
func (win *ApplicationWindow) NewestComic() {
	// Make it clear that we are checking for a new comic.
	win.header.SetTitle(l("Checking for new comic..."))
	win.ShowLoading()

	const refreshRate = time.Second
	newestComic, err := cache.CheckForNewestComicInfo(refreshRate)
	if err != nil {
		log.Print("error jumping to newest comic: ", err)
	}

	win.SetComic(newestComic.Num)
}

// RandomComic sets the current comic to a random comic.
func (win *ApplicationWindow) RandomComic() {
	today := time.Now()
	if today.Month() == time.April && today.Day() == 1 {
		win.SetComic(4) // chosen by fair dice roll.
		return          // guaranteed to be random.
	}

	win.ShowLoading()
	newestComic, _ := cache.NewestComicInfoFromCache()
	if newestComic.Num <= 0 {
		win.FirstComic()
		return
	}
	win.SetComic(rand.Intn(newestComic.Num) + 1)
}

// SetComic sets the current comic to the given comic.
func (win *ApplicationWindow) SetComic(n int) {
	win.state.ComicNumber = n

	// Make it clear that we are loading a comic.
	win.ShowLoading()
	win.header.SetSubtitle(strconv.Itoa(n))

	// Update UI to reflect new current comic.
	win.navigationBar.UpdateButtonState()
	win.bookmarksMenu.UpdateBookmarkButton()

	go func() {
		var err error

		// Add the DisplayComic function to the event loop so our UI gets
		// updated with the new comic.
		defer glib.IdleAddPriority(glib.PRIORITY_DEFAULT, win.DisplayComic)

		// Make sure we are the only ones changing win.comic.
		win.comicMutex.Lock()
		defer win.comicMutex.Unlock()

		win.comic, err = cache.ComicInfo(n)
		if err != nil {
			log.Print("error downloading comic info: ", n)
			return
		}

		_, err = os.Stat(cache.ComicImagePath(n))
		if err == nil {
			return
		}
		if !os.IsNotExist(err) {
			log.Print("error finding comic image in cache: ", err)
			return
		}

		err = cache.DownloadComicImage(n)
		if err != nil {
			log.Print("error downloading comic image: ", err)
			// We can be sneaky if we get an error, we use SafeTitle for window
			// title, but we can leave Title alone so the properties dialog can
			// still be correct.
			win.comic.SafeTitle = l("Connect to the internet to download comic image")
		}
	}()
}

// ShowLoading makes the window indicate that it is loading.
func (win *ApplicationWindow) ShowLoading() {
	win.header.SetTitle(l("Loading comic..."))
	win.comicContainer.ShowLoadingScreen()
}

// DisplayComic updates the UI to show the contents of win.comic.
func (win *ApplicationWindow) DisplayComic() {
	win.comicMutex.RLock()
	defer win.comicMutex.RUnlock()

	win.header.SetTitle(win.comic.SafeTitle)
	win.header.SetSubtitle(strconv.Itoa(win.comic.Num))
	win.comicContainer.SetTooltipText(win.comic.Alt)
	win.navigationBar.UpdateButtonState()

	// If the comic has a link, lets give the option of visiting it.
	win.actions["open-link"].SetEnabled(win.comic.Link != "")

	if win.properties != nil {
		win.properties.Update()
	}

	err := win.comicContainer.DrawComic(win.comic.Num, win.app.DarkMode())
	if err != nil {
		log.Print("error drawing comic: ", err)
	}
}

func (win *ApplicationWindow) ZoomIn() {
	win.state.ImageScale = win.comicContainer.ZoomIn()
	win.updateZoomButtonStatus()
}

func (win *ApplicationWindow) ZoomOut() {
	win.state.ImageScale = win.comicContainer.ZoomOut()
	win.updateZoomButtonStatus()
}

func (win *ApplicationWindow) ZoomReset() {
	win.state.ImageScale = win.comicContainer.SetScale(1)
	win.updateZoomButtonStatus()
}

func (win *ApplicationWindow) updateZoomButtonStatus() {
	err := win.windowMenu.zoomBox.SetCurrentZoom(win.state.ImageScale)
	if err != nil {
		log.Printf("error calling ZoomBox.SetCurrentZoom(%v): %v", win.state.ImageScale, err)
	}
	err = win.comicContainer.contextMenu.zoomBox.SetCurrentZoom(win.state.ImageScale)
	if err != nil {
		log.Printf("error calling ZoomBox.SetCurrentZoom(%v): %v", win.state.ImageScale, err)
	}
	win.actions["zoom-in"].SetEnabled(win.state.ImageScale < ImageScaleMax)
	win.actions["zoom-out"].SetEnabled(win.state.ImageScale > ImageScaleMin)
	win.actions["zoom-reset"].SetEnabled(win.state.ImageScale != 1)
}

// Explain opens a link to explainxkcd.com in the user's web browser.
func (win *ApplicationWindow) Explain() {
	win.app.OpenURL(fmt.Sprintf("https://www.explainxkcd.com/%v/", win.comicNumber()))
}

// OpenLink opens the comic's Link in the user's web browser.
func (win *ApplicationWindow) OpenLink() {
	win.comicMutex.RLock()
	link := win.comic.Link
	win.comicMutex.RUnlock()

	win.app.OpenURL(link)
}

// comicNumber returns the number of the current comic in a thread-safe way. Do
// not call this method if you already hold win.comicMutex.
func (win *ApplicationWindow) comicNumber() int {
	win.comicMutex.RLock()
	defer win.comicMutex.RUnlock()

	return win.comic.Num
}

// Dispose releases all references in the Window struct. This is needed to
// mitigate a memory leak when closing windows.
func (win *ApplicationWindow) Dispose() {
	if win == nil {
		return
	}

	win.ApplicationWindow = nil

	win.app = nil
	win.comic = nil
	win.actions = nil
	win.header = nil
	win.navigationBar.Dispose()
	win.navigationBar = nil
	win.searchMenu.Dispose()
	win.searchMenu = nil
	win.bookmarksMenu.Dispose()
	win.bookmarksMenu = nil
	win.windowMenu.Dispose()
	win.windowMenu = nil
	win.comicContainer.Dispose()
	win.comicContainer = nil
	win.properties.Dispose()
	win.properties = nil

	runtime.GC()
}
