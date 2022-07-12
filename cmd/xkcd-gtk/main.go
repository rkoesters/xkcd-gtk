// xkcd-gtk is a xkcd comic viewer app written in Go with GTK+3 (using the gotk3
// bindings).
package main

import (
	"flag"
	"fmt"
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

var (
	gtkDebug = flag.String("gtk-debug", "", "Set GTK_DEBUG environment variable to provided string.")
	version  = flag.Bool("version", false, "Print app version and exit.")
)

func usage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  %s [flags...]\n", os.Args[0])
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	rand.Seed(time.Now().Unix())
	log.Init()
	build.Init()
	paths.Init(build.AppID)

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() > 0 {
		log.Print("unexpected command line arguments: ", flag.Args())
		flag.Usage()
		os.Exit(1)
	}

	if *version {
		fmt.Printf("%v version %v\n", build.AppID, build.Version())
		os.Exit(0)
	}

	if *gtkDebug != "" {
		os.Setenv("GTK_DEBUG", *gtkDebug)
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

	os.Exit(app.Run(nil))
}
