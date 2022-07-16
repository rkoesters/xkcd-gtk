// xkcd-gtk is a xkcd comic viewer app written in Go with GTK+3 (using the gotk3
// bindings).
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/build"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"github.com/rkoesters/xkcd-gtk/internal/widget"
)

var (
	gdkDebug = flag.String("gdk-debug", "", "Behave as if the GDK_DEBUG env variable was set to the provided string.")
	gtkDebug = flag.String("gtk-debug", "", "Behave as if the GTK_DEBUG env variable was set to the provided string.")
	service  = flag.Bool("gapplication-service", false, "Start as a D-Bus service.")
	version  = flag.Bool("version", false, "Print app version and exit.")
)

func usage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "Usage of %s:\n\n", os.Args[0])
	fmt.Fprintf(w, "  %s [flags...]\n\n", os.Args[0])
	fmt.Fprintf(w, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() > 0 {
		fmt.Fprintf(flag.CommandLine.Output(), "unexpected command line arguments: %v\n", flag.Args())
		flag.Usage()
		os.Exit(1)
	}

	log.Init()
	build.Init()
	paths.Init(build.AppID())

	if *version {
		fmt.Printf("%v version %v\n", build.AppID(), build.Version())
		os.Exit(0)
	}

	if *gdkDebug != "" {
		os.Setenv("GDK_DEBUG", *gdkDebug)
	}
	if *gtkDebug != "" {
		os.Setenv("GTK_DEBUG", *gtkDebug)
	}

	glib.InitI18n(build.AppID(), paths.LocaleDir())
	glib.SetPrgname(build.AppID())
	glib.SetApplicationName(widget.AppName())
	gtk.WindowSetDefaultIconName(build.AppID())

	app, err := widget.NewApplication(build.AppID(), *service)
	if err != nil {
		log.Fatal("error creating application: ", err)
	}
	// Tell glib that this is the process's main application.
	app.SetDefault()

	os.Exit(app.Run(nil))
}
