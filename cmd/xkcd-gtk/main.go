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
	os.Exit(app.application.Run(os.Args))
}
