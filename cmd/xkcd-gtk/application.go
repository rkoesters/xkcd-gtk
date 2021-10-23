package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/bookmarks"
	"github.com/rkoesters/xkcd-gtk/internal/build"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"github.com/rkoesters/xkcd-gtk/internal/search"
	"github.com/rkoesters/xkcd-gtk/internal/settings"
	"github.com/rkoesters/xkcd-gtk/internal/style"
	"github.com/rkoesters/xkcd-gtk/internal/widget"
)

const appID = "com.github.rkoesters.xkcd-gtk"

var appName = l("Comic Sticks")

// Application holds onto our GTK representation of our application.
type Application struct {
	application *gtk.Application
	gtkSettings *gtk.Settings
	actions     map[string]*glib.SimpleAction

	aboutDialog     *gtk.AboutDialog
	shortcutsWindow *gtk.ShortcutsWindow

	settings  settings.Settings
	bookmarks bookmarks.List
}

// NewApplication creates an instance of our GTK Application.
func NewApplication() (*Application, error) {
	var app Application
	var err error

	app.application, err = gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return nil, err
	}

	// Initialize our application actions.
	app.actions = make(map[string]*glib.SimpleAction)
	registerAction := func(name string, fn interface{}) {
		action := glib.SimpleActionNew(name, nil)
		action.Connect("activate", fn)

		app.actions[name] = action
		app.application.AddAction(action)
	}

	registerAction("new-window", app.Activate)
	registerAction("open-about-xkcd", app.OpenAboutXKCD)
	registerAction("open-blog", app.OpenBlog)
	registerAction("open-store", app.OpenStore)
	registerAction("open-what-if", app.OpenWhatIf)
	registerAction("quit", app.Quit)
	registerAction("show-about", app.ShowAbout)
	registerAction("show-shortcuts", app.ShowShortcuts)
	registerAction("toggle-dark-mode", app.ToggleDarkMode)

	// Initialize our application accelerators.
	app.application.SetAccelsForAction("app.new-window", []string{"<Control>n"})
	app.application.SetAccelsForAction("app.quit", []string{"<Control>q"})
	app.application.SetAccelsForAction("app.show-shortcuts", []string{"<Control>question"})
	app.application.SetAccelsForAction("app.toggle-dark-mode", []string{"<Control>d"})
	app.application.SetAccelsForAction("win.show-properties", []string{"<Control>p"})

	// Connect startup signal to our methods.
	app.application.Connect("startup", style.InitCSS)
	app.application.Connect("startup", app.SetupAppMenu)
	app.application.Connect("startup", app.LoadSettings)
	app.application.Connect("startup", app.LoadBookmarks)
	app.application.Connect("startup", app.SetupCache)

	// Connect shutdown signal to our methods.
	app.application.Connect("shutdown", app.SaveSettings)
	app.application.Connect("shutdown", app.SaveBookmarks)
	app.application.Connect("shutdown", app.CloseCache)

	// Connect activate signal to our methods.
	app.application.Connect("activate", app.Activate)

	return &app, nil
}

// SetupAppMenu creates an AppMenu if the environment wants it.
func (app *Application) SetupAppMenu() {
	if app.application.PrefersAppMenu() {
		menu, err := widget.NewAppMenu()
		if err != nil {
			log.Fatal("error creating app menu: ", err)
		}
		app.application.SetAppMenu(menu)
	}
}

// SetupCache initializes the comic cache and the search index.
func (app *Application) SetupCache() {
	err := cache.Init(search.Index)
	if err != nil {
		log.Print("error initializing comic cache: ", err)
	}

	err = search.Init()
	if err != nil {
		log.Print("error initializing search index: ", err)
	}

	// Asynchronously fill the comic metadata cache and search index.
	search.Load(app.application)
}

// CloseCache closes the search index and comic cache.
func (app *Application) CloseCache() {
	err := search.Close()
	if err != nil {
		log.Print("error closing search index: ", err)
	}

	err = cache.Close()
	if err != nil {
		log.Print("error closing comic cache: ", err)
	}
}

// Activate creates and presents a new window to the user.
func (app *Application) Activate() {
	win, err := NewWindow(app)
	if err != nil {
		log.Print("error creating window: ", err)
		return
	}
	win.window.Present()
}

// ToggleDarkMode toggles the value of "gtk-application-prefer-dark-theme".
func (app *Application) ToggleDarkMode() {
	previous := app.DarkMode()

	// Setting 'gtk-application-prefer-dark-theme' will trigger a call to
	// win.DrawComic which will call app.DarkMode again, which will then
	// update app.settings.DarkMode (which effectively serves as a cache of
	// 'gtk-application-prefer-dark-theme').
	err := app.gtkSettings.SetProperty("gtk-application-prefer-dark-theme", !previous)
	if err != nil {
		log.Print("error setting dark mode state: ", err)
	}
}

// DarkMode returns whether the application has dark mode enabled.
func (app *Application) DarkMode() bool {
	darkMode := app.settings.DarkMode

	// Ask GTK whether it is using a dark theme.
	darkModeIface, err := app.gtkSettings.GetProperty("gtk-application-prefer-dark-theme")
	if err == nil {
		var ok bool
		darkMode, ok = darkModeIface.(bool)
		if !ok {
			log.Print("failed to interpret dark mode state")
			darkMode = app.settings.DarkMode
		}
	} else {
		log.Print("error getting dark mode state: ", err)
	}

	// Sync app.settings.DarkMode with the value of
	// 'gtk-application-prefer-dark-theme'.
	app.settings.DarkMode = darkMode

	return darkMode
}

// Quit closes all windows so the application can close.
func (app *Application) Quit() {
	// Close the active window so that it has a chance to save its state.
	win := app.application.GetActiveWindow()
	if win != nil {
		parent, _ := win.GetTransientFor()
		if parent != nil {
			win = parent
		}

		win.Close()
	}

	// Quit the application.
	glib.IdleAdd(app.application.Quit)
}

// LoadSettings tries to load our settings from disk.
func (app *Application) LoadSettings() {
	var err error

	checkForMisplacedSettings()

	// Read settings from disk.
	app.settings.ReadFile(settingsPath())

	// Get reference to Gtk's settings.
	app.gtkSettings, err = gtk.SettingsGetDefault()
	if err == nil {
		// Apply Dark Mode setting.
		err = app.gtkSettings.SetProperty("gtk-application-prefer-dark-theme", app.settings.DarkMode)
		if err != nil {
			log.Print("error setting dark mode state: ", err)
		}
	} else {
		log.Print("error querying gtk settings: ", err)
	}
}

// SaveSettings tries to save our settings to disk.
func (app *Application) SaveSettings() {
	err := paths.EnsureConfigDir()
	if err != nil {
		log.Printf("error saving settings: %v", err)
	}

	err = app.settings.WriteFile(settingsPath())
	if err != nil {
		log.Printf("error saving settings: %v", err)
	}
}

// LoadBookmarks tries to load our bookmarks from disk.
func (app *Application) LoadBookmarks() {
	checkForMisplacedBookmarks()

	app.bookmarks = bookmarks.New()
	app.bookmarks.ReadFile(bookmarksPath())
}

// SaveBookmarks tries to save our bookmarks to disk.
func (app *Application) SaveBookmarks() {
	err := paths.EnsureDataDir()
	if err != nil {
		log.Printf("error saving bookmarks: %v", err)
	}

	err = app.bookmarks.WriteFile(bookmarksPath())
	if err != nil {
		log.Printf("error saving bookmarks: %v", err)
	}
}

// ShowShortcuts shows a shortcuts window to the user.
func (app *Application) ShowShortcuts() {
	var err error
	if app.shortcutsWindow == nil {
		app.shortcutsWindow, err = widget.NewShortcutsWindow()
		if err != nil {
			log.Print("error creating shortcuts window: ", err)
			return
		}

		// We want to keep the shortcuts window around in case we want
		// to show it again.
		app.shortcutsWindow.HideOnDelete()
		app.shortcutsWindow.Connect("hide", func() {
			app.application.RemoveWindow(&app.shortcutsWindow.Window)
		})
	}

	app.application.AddWindow(&app.shortcutsWindow.Window)
	app.shortcutsWindow.Present()
}

// ShowAbout shows our application info to the user.
func (app *Application) ShowAbout() {
	var err error

	if app.aboutDialog == nil {
		app.aboutDialog, err = widget.NewAboutDialog(appID, appName, build.Version())
		if err != nil {
			log.Print("error creating about dialog: ", err)
			return
		}

		// We want to keep the about dialog around in case we want to
		// show it again.
		app.aboutDialog.HideOnDelete()
		app.aboutDialog.Connect("response", app.aboutDialog.Hide)
		app.aboutDialog.Connect("hide", func() {
			app.application.RemoveWindow(&app.aboutDialog.Window)
		})
	}

	// Set our parent window as the active window, but avoid accidentally
	// setting ourself as the parent window.
	win := app.application.GetActiveWindow()
	if win.Native() != app.aboutDialog.Native() {
		app.aboutDialog.SetTransientFor(win)
	}

	app.application.AddWindow(&app.aboutDialog.Window)
	app.aboutDialog.Present()
}

const (
	whatIfLink = "https://what-if.xkcd.com/"
	blogLink   = "https://blog.xkcd.com/"
	storeLink  = "https://store.xkcd.com/"
	aboutLink  = "https://xkcd.com/about/"
)

// OpenWhatIf opens whatifLink in the user's web browser.
func (app *Application) OpenWhatIf() {
	openURL(whatIfLink)
}

// OpenBlog opens blogLink in the user's web browser.
func (app *Application) OpenBlog() {
	openURL(blogLink)
}

// OpenStore opens storeLink in the user's web browser.
func (app *Application) OpenStore() {
	openURL(storeLink)
}

// OpenAboutXKCD opens aboutLink in the user's web browser.
func (app *Application) OpenAboutXKCD() {
	openURL(aboutLink)
}
