package main

import (
	"github.com/gotk3/gotk3/gtk"
)

type Goto struct {
	parent *Window
	dialog *gtk.Dialog
	entry  *gtk.Entry
}

func NewGoto(parent *Window) (*Goto, error) {
	var err error
	gt := new(Goto)
	gt.parent = parent

	gt.dialog, err = gtk.DialogNew()
	if err != nil {
		return nil, err
	}
	gt.dialog.SetTransientFor(parent.win)
	gt.dialog.SetDefaultSize(500, 200)

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 12)
	if err != nil {
		return nil, err
	}
	gt.entry, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	gt.entry.SetPlaceholderText("Comic #")
	gt.entry.SetHExpand(true)
	box.PackStart(gt.entry, true, false, 12)
	_, err = gt.dialog.AddButton("Cancel", 0)
	if err != nil {
		return nil, err
	}
	_, err = gt.dialog.AddButton("Go", 1)
	if err != nil {
		return nil, err
	}
	box.ShowAll()

	contentArea, err := gt.dialog.GetContentArea()
	if err != nil {
		return nil, err
	}
	contentArea.Add(box)

	return gt, nil
}
