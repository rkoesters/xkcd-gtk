// xkcd-gtk is a XKCD Comic Viewer App written in Go and GTK+3 (using
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
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Make sure our random number generator is seeded.
	rand.Seed(time.Now().Unix())

	glib.SetApplicationName("XKCD Viewer")
	gtk.WindowSetDefaultIconName("xkcd-gtk")

	// Create and run our application.
	app, err := NewApplication()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(app.GtkApp.Run(os.Args))
}
