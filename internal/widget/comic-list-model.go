package widget

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	comicListColumnNumber = iota
	comicListColumnTitle
)

type ComicListModel struct {
	*gtk.ListStore
}

func NewComicListModel() (*ComicListModel, error) {
	super, err := gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING)
	if err != nil {
		return nil, err
	}
	return &ComicListModel{
		ListStore: super,
	}, nil
}

func (clm *ComicListModel) AppendComic(comicNum int, comicTitle string) error {
	return clm.Set(
		clm.Append(),
		[]int{comicListColumnNumber, comicListColumnTitle},
		append([]any{}, comicNum, comicTitle),
	)
}
