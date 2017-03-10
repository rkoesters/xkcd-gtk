package main

//go:generate go-bindata -nomemcopy data/...

import (
	"flag"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var number = flag.Int("n", 0, "Comic number.")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	gtk.Init(nil)
	gtk.WindowSetDefaultIconName("xkcd-gtk")

	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(1)
	}

	rand.Seed(time.Now().Unix())

	viewer, err := New()
	if err != nil {
		log.Fatal(err)
	}

	viewer.win.ShowAll()

	go func() {
		viewer.SetComic(*number)
	}()

	gtk.Main()
}

func showAboutDialog() {
	abt, err := gtk.AboutDialogNew()
	if err != nil {
		log.Print(err)
		return
	}

	abt.SetProgramName("XKCD Viewer")
	abt.SetLogoIconName("xkcd-gtk")
	abt.SetVersion("0.2")
	abt.SetComments("A simple XKCD comic reader for GNOME")
	abt.SetWebsite("https://github.com/rkoesters/xkcd-gtk")
	abt.SetAuthors([]string{"Ryan Koesters"})
	abt.SetCopyright("Copyright © 2015-2017 Ryan Koesters")
	abt.SetLicenseType(gtk.LICENSE_GPL_3_0)

	abt.Show()
}
