package widget

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/rkoesters/xkcd-gtk/internal/style"
	"strconv"
)

func NewComicListButton(id int, title string, comicSetter func(int), idWidth int) (*gtk.ModelButton, error) {
	clb, err := gtk.ModelButtonNew()
	if err != nil {
		return nil, err
	}
	clb.Connect("clicked", func() { comicSetter(id) })

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		return nil, err
	}

	labelID, err := gtk.LabelNew(strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	labelID.SetXAlign(1)
	labelID.SetWidthChars(idWidth)
	box.Add(labelID)

	sep, err := gtk.SeparatorNew(gtk.ORIENTATION_VERTICAL)
	if err != nil {
		return nil, err
	}
	sep.SetMarginStart(style.PaddingComicListButton)
	sep.SetMarginEnd(style.PaddingComicListButton)
	box.Add(sep)

	labelTitle, err := gtk.LabelNew(title)
	if err != nil {
		return nil, err
	}
	labelTitle.SetEllipsize(pango.ELLIPSIZE_END)
	box.Add(labelTitle)

	child, err := clb.GetChild()
	if err != nil {
		return nil, err
	}
	clb.Remove(child)
	clb.Add(box)

	return clb, nil
}
