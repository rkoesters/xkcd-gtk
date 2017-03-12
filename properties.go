package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd"
	"strings"
	"time"
)

func NewPropertiesDialog(parent *gtk.ApplicationWindow, comic *xkcd.Comic) (*gtk.Dialog, error) {
	d, err := gtk.DialogNew()
	if err != nil {
		return nil, err
	}
	d.SetTransientFor(parent)
	d.SetDefaultSize(600, 500)

	scwin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	scwin.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)

	grid, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}
	grid.SetColumnSpacing(24)
	grid.SetRowSpacing(12)
	grid.SetMarginBottom(24)
	grid.SetMarginEnd(24)
	grid.SetMarginStart(24)
	grid.SetMarginTop(24)

	addRowToGrid(grid, 0, "Number", comic.Num)
	addRowToGrid(grid, 1, "Title", comic.Title)
	addRowToGrid(grid, 2, "Image", comic.Img)
	addRowToGrid(grid, 3, "Alt Text", comic.Alt)
	addRowToGrid(grid, 4, "Date", formatDate(comic.Year, comic.Month, comic.Day))
	addRowToGrid(grid, 5, "News", comic.News)
	addRowToGrid(grid, 6, "Link", comic.Link)
	addRowToGrid(grid, 7, "Transcript", comic.Transcript)

	scwin.Add(grid)

	box, err := d.GetContentArea()
	if err != nil {
		return nil, err
	}
	box.Add(scwin)
	// Not sure why this line is needed, but the dialog will be blank
	// without it.
	box.SetProperty("orientation", gtk.ORIENTATION_HORIZONTAL)
	box.ShowAll()

	return d, nil
}

func addRowToGrid(grid *gtk.Grid, row int, key interface{}, val interface{}) error {
	keyLabel, err := gtk.LabelNew(fmt.Sprint(key))
	if err != nil {
		return err
	}
	keyLabel.SetHAlign(gtk.ALIGN_END)
	keyLabel.SetVAlign(gtk.ALIGN_START)
	valLabel, err := gtk.LabelNew(fmt.Sprint(val))
	if err != nil {
		return err
	}
	valLabel.SetHAlign(gtk.ALIGN_START)
	valLabel.SetVAlign(gtk.ALIGN_START)
	valLabel.SetLineWrap(true)
	valLabel.SetSelectable(true)

	grid.Attach(keyLabel, 0, row, 1, 1)
	grid.Attach(valLabel, 1, row, 1, 1)

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
