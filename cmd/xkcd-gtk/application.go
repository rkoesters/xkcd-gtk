package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/skratchdot/open-golang/open"
	"log"
)

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

	actionFuncs := map[string]interface{}{
		"new-window":   app.Activate,
		"open-blog":    app.OpenBlog,
		"open-store":   app.OpenStore,
		"open-what-if": app.OpenWhatIf,
		"show-about":   app.ShowAboutDialog,
	}

	app.actions = make(map[string]*glib.SimpleAction)
	for name, function := range actionFuncs {
		action := glib.SimpleActionNew(name, nil)
		action.Connect("activate", function)

		app.actions[name] = action
		app.application.AddAction(action)
	}

	app.application.Connect("startup", app.LoadCSS)
	app.application.Connect("startup", app.LoadSearchIndex)
	app.application.Connect("activate", app.Activate)

	return &app, nil
}

// Activate creates and presents a new window to the user.
func (app *Application) Activate() {
	win, err := NewWindow(app)
	if err != nil {
		log.Fatal(err)
	}
	win.window.Present()
}

const (
	whatIfLink = "https://what-if.xkcd.com/"
	blogLink   = "https://blog.xkcd.com/"
	storeLink  = "https://store.xkcd.com/"
)

// OpenWhatIf opens whatifLink in the user's web browser.
func (app *Application) OpenWhatIf() {
	err := open.Start(whatIfLink)
	if err != nil {
		log.Print(err)
	}
}

// OpenBlog opens blogLink in the user's web browser.
func (app *Application) OpenBlog() {
	err := open.Start(blogLink)
	if err != nil {
		log.Print(err)
	}
}

// OpenStore opens storeLink in the user's web browser.
func (app *Application) OpenStore() {
	err := open.Start(storeLink)
	if err != nil {
		log.Print(err)
	}
}
