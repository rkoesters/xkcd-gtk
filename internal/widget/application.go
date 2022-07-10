package widget

import (
	"errors"
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
	"os"
)

// AppName is the user-visible name of this application.
func AppName() string { return l("Comic Sticks") }

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

	app.application, err = gtk.ApplicationNew(build.AppID, glib.APPLICATION_FLAGS_NONE)
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

	// Connect application signal handlers.
	app.application.Connect("startup", app.Startup)
	app.application.Connect("shutdown", app.Shutdown)
	app.application.Connect("activate", app.Activate)

	return &app, nil
}

// Startup is called when the "startup" signal is emitted.
func (app *Application) Startup() {
	err := app.SetupAppMenu()
	if err != nil {
		log.Fatal("error creating app menu: ", err)
	}
	app.LoadSettings()
	style.InitCSS(app.DarkMode())
	app.gtkSettings.Connect("notify::gtk-application-prefer-dark-theme", app.DarkModeChanged)
	app.LoadBookmarks()
	app.SetupCache()
}

// Shutdown is called when the "shutdown" signal is emitted.
func (app *Application) Shutdown() {
	app.SaveSettings()
	app.SaveBookmarks()
	app.CloseCache()
}

// SetDefault is a wrapper around glib.Application.SetDefault().
func (app *Application) SetDefault() { app.application.SetDefault() }

// Run is a wrapper around glib.Application.Run().
func (app *Application) Run(args []string) int {
	return app.application.Run(args)
}

// PrefersAppMenu is a wrapper around gtk.Application.PrefersAppMenu().
func (app *Application) PrefersAppMenu() bool {
	return app.application.PrefersAppMenu() || build.Options["always-prefer-app-menu"] == "true"
}

// SetupAppMenu creates an AppMenu if the environment wants it.
func (app *Application) SetupAppMenu() error {
	if !app.PrefersAppMenu() {
		return nil
	}
	menu, err := NewAppMenu()
	if err != nil {
		return err
	}
	app.application.SetAppMenu(menu)
	return nil
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
	err = search.Load(app.application)
	if err != nil {
		log.Print("error building search index: ", err)
	}
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
	win, err := NewApplicationWindow(app)
	if err != nil {
		log.Print("error creating window: ", err)
		return
	}
	win.window.Present()
}

// DarkModeChanged is called when gtk-application-prefer-dark-theme is changed.
func (app *Application) DarkModeChanged() {
	darkMode := app.DarkMode()
	log.Debugf("DarkModeChanged() -> %v", darkMode)
	err := style.UpdateCSS(darkMode)
	if err != nil {
		log.Printf("error calling style.UpdateCSS(darkMode=%v) -> %v", darkMode, err)
	}
}

// ToggleDarkMode toggles the value of "gtk-application-prefer-dark-theme".
func (app *Application) ToggleDarkMode() {
	app.SetDarkMode(!app.DarkMode())
}

// SetDarkMode sets the value of "gtk-application-prefer-dark-theme" to the
// darkMode argument.
func (app *Application) SetDarkMode(darkMode bool) {
	log.Debugf("SetDarkMode(darkMode=%v)", darkMode)
	// Change the dark mode setting in one of the next iterations of the
	// event loop (i.e. do not block) so that the style does not change in
	// the middle of any ongoing animations (e.g. a switch toggling or a
	// menu closing).
	go glib.IdleAdd(func() {
		log.Debugf("SetDarkMode(darkMode=%v).func()", darkMode)
		// Setting 'gtk-application-prefer-dark-theme' will trigger a
		// call to win.DrawComic which will call app.DarkMode again,
		// which will then update app.settings.DarkMode (which
		// effectively serves as a cache of
		// 'gtk-application-prefer-dark-theme').
		err := app.gtkSettings.SetProperty("gtk-application-prefer-dark-theme", darkMode)
		if err != nil {
			log.Print("error setting dark mode state: ", err)
		}
	})
}

// DarkMode returns whether the application has dark mode enabled.
func (app *Application) DarkMode() bool {
	// Ask GTK whether it is using a dark theme.
	darkModeIface, err := app.gtkSettings.GetProperty("gtk-application-prefer-dark-theme")
	if err != nil {
		log.Print("error getting dark mode state: ", err)
		return app.settings.DarkMode
	}

	darkMode, ok := darkModeIface.(bool)
	if !ok {
		log.Print("failed to interpret dark mode state")
		return app.settings.DarkMode
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

	paths.CheckForMisplacedSettings()

	// Read settings from disk.
	err = app.settings.ReadFile(paths.Settings())
	if err != nil {
		log.Print("error reading app settings: ", err)
	}

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

	err = app.settings.WriteFile(paths.Settings())
	if err != nil {
		log.Printf("error saving settings: %v", err)
	}
}

// LoadBookmarks tries to load our bookmarks from disk.
func (app *Application) LoadBookmarks() {
	paths.CheckForMisplacedBookmarks()

	app.bookmarks = bookmarks.New()
	err := app.bookmarks.ReadFile(paths.Bookmarks())
	if err != nil {
		log.Print("error reading bookmarks: ", err)
	}
}

// SaveBookmarks tries to save our bookmarks to disk.
func (app *Application) SaveBookmarks() {
	err := paths.EnsureDataDir()
	if err != nil {
		log.Printf("error saving bookmarks: %v", err)
	}

	err = app.bookmarks.WriteFile(paths.Bookmarks())
	if err != nil {
		log.Printf("error saving bookmarks: %v", err)
	}
}

// ShowShortcuts shows a shortcuts window to the user.
func (app *Application) ShowShortcuts() {
	var err error
	if app.shortcutsWindow == nil {
		app.shortcutsWindow, err = NewShortcutsWindow(app.application.RemoveWindow)
		if err != nil {
			log.Print("error creating shortcuts window: ", err)
			return
		}
	}

	app.application.AddWindow(app.shortcutsWindow)
	app.shortcutsWindow.Present()
}

// ShowAbout shows our application info to the user.
func (app *Application) ShowAbout() {
	var err error

	if app.aboutDialog == nil {
		app.aboutDialog, err = NewAboutDialog(app.application.RemoveWindow)
		if err != nil {
			log.Print("error creating about dialog: ", err)
			return
		}
	}

	// Set our parent window as the active window, but avoid accidentally
	// setting ourself as the parent window.
	win := app.application.GetActiveWindow()
	if win.Native() != app.aboutDialog.Native() {
		app.aboutDialog.SetTransientFor(win)
	}

	app.application.AddWindow(app.aboutDialog)
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

func (app *Application) gtkTheme() (string, error) {
	theme := os.Getenv("GTK_THEME")
	if theme != "" {
		return theme, nil
	}
	themeIface, err := app.gtkSettings.GetProperty("gtk-theme-name")
	if err != nil {
		return "", err
	}
	var ok bool
	theme, ok = themeIface.(string)
	if !ok {
		return "", errors.New("error converting gtk-theme-name to a string")
	}
	return theme, nil
}
