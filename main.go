package main

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

	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(1)
	}

	if *number == 0 {
		c, err := getNewestComicInfo()
		if err != nil {
			log.Fatal(err)
		}
		*number = c.Num
	}

	rand.Seed(time.Now().Unix())

	viewer, err := New()
	if err != nil {
		log.Fatal(err)
	}

	viewer.SetComic(*number)
	viewer.win.ShowAll()

	gtk.Main()
}

func showAboutDialog() {
	builder, err := gtk.BuilderNew()
	if err != nil {
		log.Print(err)
		return
	}
	err = builder.AddFromFile("about.ui")
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
