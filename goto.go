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
	gt.dialog.SetResizable(false)

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 12)
	if err != nil {
		return nil, err
	}
	box.SetMarginStart(6)
	box.SetMarginEnd(6)
	label, err := gtk.LabelNew("Go to comic number")
	if err != nil {
		return nil, err
	}
	box.Add(label)
	gt.entry, err = gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	gt.entry.SetActivatesDefault(true)
	gt.entry.SetPlaceholderText("Comic #")
	box.Add(gt.entry)
	_, err = gt.dialog.AddButton("Cancel", 0)
	if err != nil {
		return nil, err
	}
	submit, err := gt.dialog.AddButton("Go", 1)
	if err != nil {
		return nil, err
	}
	submitStyle, err := submit.GetStyleContext()
	if err != nil {
		return nil, err
	}
	submitStyle.AddClass("suggested-action")
	submit.SetCanDefault(true)
	submit.GrabDefault()
	box.ShowAll()

	contentArea, err := gt.dialog.GetContentArea()
	if err != nil {
		return nil, err
	}
	contentArea.Add(box)

	return gt, nil
}
