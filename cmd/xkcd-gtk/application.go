package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/bookmarks"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/search"
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
	darkModeIface, err := app.gtkSettings.GetProperty("gtk-application-prefer-dark-theme")
	if err != nil {
		log.Print("error getting dark mode state: ", err)
		return
	}

	darkMode, ok := darkModeIface.(bool)
	if !ok {
		log.Print("failed to interpret dark mode state")
		return
	}

	err = app.gtkSettings.SetProperty("gtk-application-prefer-dark-theme", !darkMode)
	if err != nil {
		log.Print("error setting dark mode state: ", err)
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
