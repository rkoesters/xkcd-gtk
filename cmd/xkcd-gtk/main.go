// xkcd-gtk is a xkcd comic viewer app written in Go with GTK+3 (using the gotk3
// bindings).
package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/build"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"math/rand"
	"os"
	"time"
)

func main() {
	// Initialize internal log package.
	log.Init()

	// Make sure the random number generator is seeded.
	rand.Seed(time.Now().Unix())

	// Initialize compile-time information provided by the build package.
	build.Parse()

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

	// Show gtk's interactive debugging window if this is a debugging build.
	args := os.Args
	if build.Debug() {
		// Insert --gtk-debug=interactive as first flag after args[0].
		args = append(args, "")
		copy(args[2:], args[1:])
		args[1] = "--gtk-debug=interactive"
	}

	// Run the event loop.
	os.Exit(app.application.Run(args))
}
