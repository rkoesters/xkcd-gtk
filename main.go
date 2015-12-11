package main

//go:generate go-bindata -nomemcopy data/...

import (
	"flag"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"math/rand"
	"os"
	"time"
)

var number = flag.Int("n", 0, "Comic number.")

func main() {
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
	builder, err := gtk.BuilderNew()
	if err != nil {
		log.Print(err)
		return
	}

	data, err := Asset("data/about.ui")
	if err != nil {
		log.Print(err)
		return
	}

	err = builder.AddFromString(string(data))
	if err != nil {
		log.Print(err)
		return
	}

	obj, err := builder.GetObject("about-dialog")
	if err != nil {
		log.Print(err)
		return
	}
	win, ok := obj.(*gtk.AboutDialog)
	if !ok {
		log.Print("error getting about-dialog")
		return
	}
	win.Show()
}
