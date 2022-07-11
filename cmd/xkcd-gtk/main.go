// xkcd-gtk is a xkcd comic viewer app written in Go with GTK+3 (using the gotk3
// bindings).
package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/build"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"github.com/rkoesters/xkcd-gtk/internal/widget"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	log.Init()
	build.Init()
	paths.Init(build.AppID)

	if build.Debug() {
		// Show gtk's interactive debugging window.
		os.Setenv("GTK_DEBUG", "interactive")
	}

	glib.InitI18n(build.AppID, paths.LocaleDir())
	glib.SetPrgname(build.AppID)
	glib.SetApplicationName(widget.AppName())
	gtk.WindowSetDefaultIconName(build.AppID)

	app, err := widget.NewApplication(build.AppID)
	if err != nil {
		log.Fatal("error creating application: ", err)
	}
	// Tell glib that this is the process's main application.
	app.SetDefault()

	os.Exit(app.Run(os.Args))
}
