package widget

import (
	"strconv"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

func NewComicListButton(id int, title string, comicSetter func(int), idWidth int) (*gtk.ModelButton, error) {
	clb, err := gtk.ModelButtonNew()
	if err != nil {
		return nil, err
	}
	clb.Connect("clicked", func() { comicSetter(id) })

	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, style.PaddingComicListButton)
	if err != nil {
		return nil, err
	}

	labelID, err := gtk.LabelNew(strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	labelID.SetXAlign(1) // align end
	labelID.SetWidthChars(idWidth)
	box.PackStart(labelID, false, false, 0)

	labelTitle, err := gtk.LabelNew(title)
	if err != nil {
		return nil, err
	}
	labelTitle.SetXAlign(0) // align start
	labelTitle.SetEllipsize(pango.ELLIPSIZE_END)
	box.PackStart(labelTitle, true, true, 0)

	child, err := clb.GetChild()
	if err != nil {
		return nil, err
	}
	clb.Remove(child)
	clb.Add(box)

	return clb, nil
}
