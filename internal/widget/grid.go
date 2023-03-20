package widget

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type Grid struct {
	*gtk.Grid
	row int
}

var _ Widget = &Grid{}

func NewGrid() (*Grid, error) {
	super, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}
	g := &Grid{
		Grid: super,
	}

	g.SetColumnSpacing(style.PaddingAuxiliaryWindow)
	g.SetRowSpacing(style.PaddingAuxiliaryWindow)
	g.SetMarginTop(style.PaddingAuxiliaryWindow)
	g.SetMarginBottom(style.PaddingAuxiliaryWindow)
	g.SetMarginStart(style.PaddingAuxiliaryWindow)
	g.SetMarginEnd(style.PaddingAuxiliaryWindow)

	return g, nil
}

func (g *Grid) AddRowToGrid(label string) (*gtk.Label, error) {
	keyLabel, err := gtk.LabelNew(label)
	if err != nil {
		return nil, err
	}
	keyLabel.SetHAlign(gtk.ALIGN_END)
	keyLabel.SetVAlign(gtk.ALIGN_START)
	valLabel, err := gtk.LabelNew("")
	if err != nil {
		return nil, err
	}
	valLabel.SetXAlign(0)
	valLabel.SetHAlign(gtk.ALIGN_START)
	valLabel.SetVAlign(gtk.ALIGN_START)
	valLabel.SetLineWrap(true)
	valLabel.SetSelectable(true)
	valLabel.SetCanFocus(false)

	g.Attach(keyLabel, 0, g.row, 1, 1)
	g.Attach(valLabel, 1, g.row, 1, 1)

	g.row++

	return valLabel, nil
}

func (g *Grid) Dispose() {
	if g == nil {
		return
	}
	g.Grid = nil
}
