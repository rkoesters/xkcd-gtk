package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Window is the main application window.
type Window struct {
	comic    *xkcd.Comic
	win      *gtk.ApplicationWindow
	hdr      *gtk.HeaderBar
	previous *gtk.Button
	next     *gtk.Button
	img      *gtk.Image
}

// New creates a new XKCD viewer window.
func NewWindow(app *Application) (*Window, error) {
	var err error

	w := new(Window)

	w.win, err = gtk.ApplicationWindowNew(app.GtkApp)
	if err != nil {
		return nil, err
	}
	w.win.SetDefaultSize(1000, 800)

	// Create HeaderBar
	w.hdr, err = gtk.HeaderBarNew()
	if err != nil {
		return nil, err
	}
	w.hdr.SetTitle("XKCD Viewer")
	w.hdr.SetShowCloseButton(true)

	w.previous, err = gtk.ButtonNewFromIconName("go-previous-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
	if err != nil {
		return nil, err
	}
	w.previous.Connect("clicked", w.PreviousComic)
	w.hdr.PackStart(w.previous)

	w.next, err = gtk.ButtonNewFromIconName("go-next-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
	if err != nil {
		return nil, err
	}
	w.next.Connect("clicked", w.NextComic)
	w.hdr.PackStart(w.next)

	randBtn, err := gtk.ButtonNewFromIconName("media-playlist-shuffle-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
	if err != nil {
		return nil, err
	}
	randBtn.Connect("clicked", w.RandomComic)
	w.hdr.PackStart(randBtn)

	w.hdr.ShowAll()
	w.win.SetTitlebar(w.hdr)

	// Create main part of window.
	scwin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}

	w.img, err = gtk.ImageNewFromIconName("emblem-synchronizing-symbolic", gtk.ICON_SIZE_DIALOG)
	if err != nil {
		return nil, err
	}
	scwin.Add(w.img)
	scwin.ShowAll()
	w.win.Add(scwin)

	return w, nil
}

// PreviousComic sets the current comic to the previous comic.
func (w *Window) PreviousComic() {
	err := w.SetComic(w.comic.Num - 1)
	if err != nil {
		log.Print(err)
	}
}

// NextComic sets the current comic to the next comic.
func (w *Window) NextComic() {
	err := w.SetComic(w.comic.Num + 1)
	if err != nil {
		log.Print(err)
	}
}

// RandomComic sets the current comic to a random comic.
func (w *Window) RandomComic() {
	c, err := getNewestComicInfo()
	if err != nil {
		log.Print(err)
		return
	}
	err = w.SetComic(rand.Intn(c.Num) + 1)
	if err != nil {
		log.Print(err)
	}
}

// ShowProperties shows a properties dialog containing all the
// information on the current comic.
func (v *Window) ShowProperties() {
	builder, err := gtk.BuilderNew()
	if err != nil {
		log.Print(err)
		return
	}
	data, err := Asset("data/properties.ui")
	if err != nil {
		log.Print(err)
		return
	}
	err = builder.AddFromString(string(data))
	if err != nil {
		log.Print(err)
		return
	}

	obj, err := builder.GetObject("properties-dialog")
	if err != nil {
		log.Print(err)
		return
	}
	dialog, ok := obj.(*gtk.Dialog)
	if !ok {
		log.Print("error getting properties-dialog")
		return
	}
	dialog.SetTransientFor(v.win)
	dialog.SetModal(true)

	number, err := getLabel(builder, "properties-number")
	if err != nil {
		log.Print(err)
		return
	}
	number.SetText(strconv.Itoa(v.comic.Num))
	title, err := getLabel(builder, "properties-title")
	if err != nil {
		log.Print(err)
		return
	}
	title.SetText(v.comic.Title)
	image, err := getLabel(builder, "properties-image")
	if err != nil {
		log.Print(err)
		return
	}
	fmtImage := fmt.Sprintf("<a href=\"%v\">%[1]v</a>", v.comic.Img)
	image.SetMarkup(fmtImage)
	alt, err := getLabel(builder, "properties-alt")
	if err != nil {
		log.Print(err)
		return
	}
	alt.SetText(v.comic.Alt)
	date, err := getLabel(builder, "properties-date")
	if err != nil {
		log.Print(err)
		return
	}
	date.SetText(formatDate(v.comic.Year, v.comic.Month, v.comic.Day))
	news, err := getLabel(builder, "properties-news")
	if err != nil {
		log.Print(err)
		return
	}
	news.SetText(v.comic.News)
	link, err := getLabel(builder, "properties-link")
	if err != nil {
		log.Print(err)
		return
	}
	if v.comic.Link != "" {
		fmtLink := fmt.Sprintf("<a href=\"%v\">%[1]v</a>", v.comic.Link)
		link.SetMarkup(fmtLink)
	}
	transcript, err := getLabel(builder, "properties-transcript")
	if err != nil {
		log.Print(err)
		return
	}
	transcript.SetText(strings.Replace(v.comic.Transcript, "\\n", "\n", -1))

	dialog.Show()
}

// SetComic sets the current comic to the given comic.
func (w *Window) SetComic(n int) error {
	var c *xkcd.Comic
	var err error
	if n == 0 {
		c, err = getNewestComicInfo()
		if err != nil {
			return err
		}
	} else {
		c, err = getComicInfo(n)
		if err != nil {
			return err
		}
	}
	w.comic = c

	imgPath, err := getComicImage(w.comic.Num)
	if err != nil {
		log.Printf("error downloading comic: %v", w.comic.Num)
	}
	w.hdr.SetSubtitle(fmt.Sprintf("#%v: %v", w.comic.Num, w.comic.Title))
	w.img.SetFromFile(imgPath)
	w.img.SetTooltipText(w.comic.Alt)

	// Enable/disable previous button.
	if w.comic.Num > 1 {
		w.previous.SetSensitive(true)
	} else {
		w.previous.SetSensitive(false)
	}

	// Enable/disable next button.
	newest, err := getNewestComicInfo()
	if err != nil {
		return err
	}
	if w.comic.Num < newest.Num {
		w.next.SetSensitive(true)
	} else {
		w.next.SetSensitive(false)
	}

	return nil
}

func getLabel(b *gtk.Builder, id string) (*gtk.Label, error) {
	obj, err := b.GetObject(id)
	if err != nil {
		return nil, err
	}
	label, ok := obj.(*gtk.Label)
	if !ok {
		return nil, fmt.Errorf("error getting label: %v", id)
	}
	return label, nil
}

// formatDate takes a year, month, and date as strings and turns them
// into a pretty date.
func formatDate(year, month, day string) string {
	date := strings.Join([]string{year, month, day}, "-")
	t, err := time.Parse("2006-1-2", date)
	if err != nil {
		return ""
	}
	return t.Format("Jan _2, 2006")
}
