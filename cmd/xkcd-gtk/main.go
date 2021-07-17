// xkcd-gtk is a xkcd comic viewer app written in Go with GTK+3 (using the gotk3
// bindings).
package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	// Make log messages include date, time, filename, and line number.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Make sure the random number generator is seeded.
	rand.Seed(time.Now().Unix())

	// Initialize the paths under which we will store app files.
	paths.Init(appID)

	// Let glib and gtk know who we are.
	glib.SetApplicationName(appName)
	gtk.WindowSetDefaultIconName(appID)

	// Create the application.
	app, err := NewApplication()
	if err != nil {
		log.Fatal("error creating application: ", err)
	}
	// Tell glib that this is the process's main application.
	app.application.SetDefault()

	// Run the event loop.
	os.Exit(app.application.Run(os.Args))
}
