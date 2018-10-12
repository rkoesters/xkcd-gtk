// xkcd-gtk is a xkcd comic viewer app written in Go with GTK+3 (using
// the gotk3 bindings).
package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	// Tell the go runtime to use as many CPUs as are available.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Make sure the random number generator is seeded.
	rand.Seed(time.Now().Unix())

	// Initialize the comic cache.
	err := initComicCache()
	if err != nil {
		log.Fatalf("failed to initialize comic cache: %v", err)
	}

	// Let glib and gtk know who we are.
	glib.SetApplicationName(appName)
	gtk.WindowSetDefaultIconName(appID)

	// Create the application.
	app, err := NewApplication()
	if err != nil {
		log.Fatal(err)
	}
	// Tell glib that this is the process's main application.
	app.application.SetDefault()

	// Run the event loop.
	ret := app.application.Run(os.Args)

	// Close the comic cache.
	closeComicCache()

	// Exit with the status code given by gtk.
	os.Exit(ret)
}
