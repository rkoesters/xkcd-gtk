package main

import (
	"errors"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Viewer is a struct holding a gtk window for viewing XKCD comics.
type Viewer struct {
	comic    *xkcd.Comic
	win      *gtk.Window
	hdr      *gtk.HeaderBar
	previous *gtk.Button
	next     *gtk.Button
	img      *gtk.Image
}

// New creates a new XKCD viewer window.
func New() (*Viewer, error) {
	v := new(Viewer)

	// Builder the gtk interface using gtk.Builder.
	builder, err := gtk.BuilderNew()
	if err != nil {
		return nil, err
	}
	data, err := Asset("data/viewer.ui")
	if err != nil {
		return nil, err
	}
	err = builder.AddFromString(string(data))
	if err != nil {
		return nil, err
	}

	// Connect the gtk signals to our functions.
	builder.ConnectSignals(map[string]interface{}{
		"PreviousComic":   v.PreviousComic,
		"NextComic":       v.NextComic,
		"RandomComic":     v.RandomComic,
		"ShowProperties":  v.ShowProperties,
		"showAboutDialog": showAboutDialog,
	})

	// We want access to Window, HeaderBar, and Image in the future,
	// so lets get access to them now.
	var ok bool
	obj, err := builder.GetObject("viewer-window")
	if err != nil {
		return nil, err
	}
	v.win, ok = obj.(*gtk.Window)
	if !ok {
		return nil, errors.New("error getting viewer-window")
	}
	obj, err = builder.GetObject("header")
	if err != nil {
		return nil, err
	}
	v.hdr, ok = obj.(*gtk.HeaderBar)
	if !ok {
		return nil, errors.New("error getting header")
	}
	obj, err = builder.GetObject("previous")
	if err != nil {
		return nil, err
	}
	v.previous, ok = obj.(*gtk.Button)
	if !ok {
		return nil, errors.New("error getting previous")
	}
	obj, err = builder.GetObject("next")
	if err != nil {
		return nil, err
	}
	v.next, ok = obj.(*gtk.Button)
	if !ok {
		return nil, errors.New("error getting next")
	}
	obj, err = builder.GetObject("comic-image")
	if err != nil {
		return nil, err
	}
	v.img, ok = obj.(*gtk.Image)
	if !ok {
		return nil, errors.New("error getting comic-image")
	}

	// Closing the window should exit the program.
	v.win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	return v, nil
}

// PreviousComic sets the current comic to the previous comic.
func (v *Viewer) PreviousComic() {
	err := v.SetComic(v.comic.Num - 1)
	if err != nil {
		log.Print(err)
	}
}

// NextComic sets the current comic to the next comic.
func (v *Viewer) NextComic() {
	err := v.SetComic(v.comic.Num + 1)
	if err != nil {
		log.Print(err)
	}
}

// RandomComic sets the current comic to a random comic.
func (v *Viewer) RandomComic() {
	c, err := getNewestComicInfo()
	if err != nil {
		log.Print(err)
		return
	}
	err = v.SetComic(rand.Intn(c.Num) + 1)
	if err != nil {
		log.Print(err)
	}
}

func (v *Viewer) ShowProperties() {
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
func (v *Viewer) SetComic(n int) error {
	c, err := getComicInfo(n)
	if err != nil {
		return err
	}
	v.comic = c

	imgPath, err := getComicImage(n)
	if err != nil {
		log.Printf("error downloading comic: %v", n)
	}
	v.hdr.SetSubtitle(fmt.Sprintf("#%v: %v", v.comic.Num, v.comic.Title))
	v.img.SetFromFile(imgPath)
	v.img.SetTooltipText(v.comic.Alt)

	// Enable/disable previous button.
	if v.comic.Num > 1 {
		v.previous.SetSensitive(true)
	} else {
		v.previous.SetSensitive(false)
	}

	// Enable/disable next button.
	newest, err := getNewestComicInfo()
	if err != nil {
		return err
	}
	if v.comic.Num < newest.Num {
		v.next.SetSensitive(true)
	} else {
		v.next.SetSensitive(false)
	}

	return nil
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
