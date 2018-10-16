package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xdg"
	"log"
)

const (
	appID   = "com.github.rkoesters.xkcd-gtk"
	appName = "Comic Sticks"
)

var appVersion = "undefined"

// Application holds onto our GTK representation of our application.
type Application struct {
	application *gtk.Application

	actions map[string]*glib.SimpleAction
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
		"new-window":     app.Activate,
		"open-blog":      app.OpenBlog,
		"open-store":     app.OpenStore,
		"open-what-if":   app.OpenWhatIf,
		"quit":           app.Quit,
		"show-about":     app.ShowAbout,
		"show-shortcuts": app.ShowShortcuts,
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
	app.application.SetAccelsForAction("app.show-shortcuts", []string{"<Control>question"})
	app.application.SetAccelsForAction("app.quit", []string{"<Control>q"})

	// Connect application signals to our methods.
	app.application.Connect("startup", app.LoadCSS)
	app.application.Connect("startup", app.SetupAppMenu)
	app.application.Connect("startup", app.SetupCache)
	app.application.Connect("shutdown", app.CloseCache)
	app.application.Connect("activate", app.Activate)

	return &app, nil
}

// SetupAppMenu creates an AppMenu if the environment wants it.
func (app *Application) SetupAppMenu() {
	if app.application.PrefersAppMenu() {
		menuSection1 := glib.MenuNew()
		menuSection1.Append("New Window", "app.new-window")

		menuSection2 := glib.MenuNew()
		menuSection2.Append("what if?", "app.open-what-if")
		menuSection2.Append("xkcd blog", "app.open-blog")
		menuSection2.Append("xkcd store", "app.open-store")

		menuSection3 := glib.MenuNew()
		menuSection3.Append("Keyboard Shortcuts", "app.show-shortcuts")
		menuSection3.Append("About "+appName, "app.show-about")
		menuSection3.Append("Quit", "app.quit")

		menu := glib.MenuNew()
		menu.AppendSectionWithoutLabel(&menuSection1.MenuModel)
		menu.AppendSectionWithoutLabel(&menuSection2.MenuModel)
		menu.AppendSectionWithoutLabel(&menuSection3.MenuModel)

		app.application.SetAppMenu(&menu.MenuModel)
	}
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
