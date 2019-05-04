package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xdg"
	"log"
)

const appID = "com.github.rkoesters.xkcd-gtk"

var (
	appName    = l("Comic Sticks")
	appVersion = "undefined"
)

// Application holds onto our GTK representation of our application.
type Application struct {
	application *gtk.Application
	gtkSettings *gtk.Settings
	actions     map[string]*glib.SimpleAction

	settings  Settings
	bookmarks Bookmarks
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
	actionFuncs := map[string]interface{}{
		"new-window":       app.Activate,
		"open-blog":        app.OpenBlog,
		"open-store":       app.OpenStore,
		"open-what-if":     app.OpenWhatIf,
		"open-about-xkcd":  app.OpenAboutXKCD,
		"quit":             app.Quit,
		"show-about":       app.ShowAbout,
		"show-shortcuts":   app.ShowShortcuts,
		"toggle-dark-mode": app.ToggleDarkMode,
	}

	app.actions = make(map[string]*glib.SimpleAction)
	for name, function := range actionFuncs {
		action := glib.SimpleActionNew(name, nil)
		action.Connect("activate", function)

		app.actions[name] = action
		app.application.AddAction(action)
	}

	// Initialize our application accelerators.
	app.application.SetAccelsForAction("app.new-window", []string{"<Control>n"})
	app.application.SetAccelsForAction("app.quit", []string{"<Control>q"})
	app.application.SetAccelsForAction("app.show-shortcuts", []string{"<Control>question"})
	app.application.SetAccelsForAction("app.toggle-dark-mode", []string{"<Control>d"})

	// Connect startup signal to our methods.
	app.application.Connect("startup", app.LoadCSS)
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

// SetupCache initializes the comic cache and the search index.
func (app *Application) SetupCache() {
	err := initComicCache()
	if err != nil {
		log.Print(err)
	}

	err = initSearchIndex()
	if err != nil {
		log.Print(err)
	}

	app.LoadSearchIndex()
}

// CloseCache closes the search index and comic cache.
func (app *Application) CloseCache() {
	err := closeSearchIndex()
	if err != nil {
		log.Print(err)
	}

	err = closeComicCache()
	if err != nil {
		log.Print(err)
	}
}

// Activate creates and presents a new window to the user.
func (app *Application) Activate() {
	win, err := NewWindow(app)
	if err != nil {
		log.Print(err)
		return
	}
	win.window.Present()
}

// ToggleDarkMode toggles the value of "gtk-application-prefer-dark-theme".
func (app *Application) ToggleDarkMode() {
	darkModeIface, err := app.gtkSettings.GetProperty("gtk-application-prefer-dark-theme")
	if err != nil {
		log.Print(err)
		return
	}

	darkMode, ok := darkModeIface.(bool)
	if !ok {
		log.Print("failed to convert darkModeIface to bool")
		return
	}

	err = app.gtkSettings.SetProperty("gtk-application-prefer-dark-theme", !darkMode)
	if err != nil {
		log.Print(err)
		return
	}
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

const (
	whatIfLink = "https://what-if.xkcd.com/"
	blogLink   = "https://blog.xkcd.com/"
	storeLink  = "https://store.xkcd.com/"
	aboutLink  = "https://xkcd.com/about/"
)

// OpenWhatIf opens whatifLink in the user's web browser.
func (app *Application) OpenWhatIf() {
	err := xdg.Open(whatIfLink)
	if err != nil {
		log.Print(err)
	}
}

// OpenBlog opens blogLink in the user's web browser.
func (app *Application) OpenBlog() {
	err := xdg.Open(blogLink)
	if err != nil {
		log.Print(err)
	}
}

// OpenStore opens storeLink in the user's web browser.
func (app *Application) OpenStore() {
	err := xdg.Open(storeLink)
	if err != nil {
		log.Print(err)
	}
}

// OpenAboutXKCD opens aboutLink in the user's web browser.
func (app *Application) OpenAboutXKCD() {
	err := xdg.Open(aboutLink)
	if err != nil {
		log.Print(err)
	}
}
