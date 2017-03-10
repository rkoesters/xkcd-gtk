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
