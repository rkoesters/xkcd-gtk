package widget

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/log"
)

type ComicListView struct {
	*gtk.TreeView

	numberColumn *gtk.TreeViewColumn
	titleColumn  *gtk.TreeViewColumn

	setComic func(int) // win.SetComic
}

var _ Widget = &ComicListView{}

func NewComicListView(comicSetter func(int)) (*ComicListView, error) {
	super, err := gtk.TreeViewNew()
	if err != nil {
		return nil, err
	}
	clv := &ComicListView{
		TreeView: super,
		setComic: comicSetter,
	}

	clv.SetHeadersVisible(false)
	clv.SetEnableSearch(false)
	clv.SetActivateOnSingleClick(true)
	clv.SetHoverSelection(true)

	const (
		xpad = 3
		ypad = 4
	)

	numberRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return nil, err
	}
	numberRenderer.SetAlignment(1, 0) // xalign right, yalign top
	numberRenderer.SetProperty("xpad", xpad)
	numberRenderer.SetProperty("ypad", ypad)
	clv.numberColumn, err = gtk.TreeViewColumnNewWithAttribute("number", numberRenderer, "text", comicListColumnNumber)
	if err != nil {
		return nil, err
	}
	clv.numberColumn.SetVisible(true)
	clv.InsertColumn(clv.numberColumn, comicListColumnNumber)

	titleRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return nil, err
	}
	titleRenderer.SetAlignment(0, 0) // xalign left, yalign top
	titleRenderer.SetProperty("xpad", xpad)
	titleRenderer.SetProperty("ypad", ypad)
	clv.titleColumn, err = gtk.TreeViewColumnNewWithAttribute("title", titleRenderer, "text", comicListColumnTitle)
	if err != nil {
		return nil, err
	}
	clv.titleColumn.SetVisible(true)
	clv.InsertColumn(clv.titleColumn, comicListColumnTitle)

	clv.Show()

	clv.Connect("row-activated", clv.rowActivated)

	return clv, nil
}

func (clv *ComicListView) Dispose() {
	if clv == nil {
		return
	}

	clv.TreeView = nil
	clv.numberColumn = nil
	clv.titleColumn = nil
	clv.setComic = nil
}

func (clv *ComicListView) SetModel(m gtk.ITreeModel) {
	clv.TreeView.SetModel(m)
	clv.numberColumn.QueueResize()
	clv.titleColumn.QueueResize()
}

func (clv *ComicListView) rowActivated(tv *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
	itm, err := clv.GetModel()
	if err != nil {
		log.Print(err)
		return
	}
	tm := itm.ToTreeModel()
	iter, err := tm.GetIter(path)
	if err != nil {
		log.Print(err)
		return
	}
	val, err := tm.GetValue(iter, comicListColumnNumber)
	if err != nil {
		log.Print(err)
		return
	}
	id, err := val.GoValue()
	if err != nil {
		log.Print(err)
		return
	}
	n, ok := id.(int)
	if !ok {
		log.Print("error converting val to int")
		return
	}
	clv.setComic(n)
}

func NewComicListScroller() (*gtk.ScrolledWindow, error) {
	scroller, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	scroller.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	scroller.SetPropagateNaturalHeight(true)
	scroller.SetMaxContentHeight(350)
	scroller.SetShadowType(gtk.SHADOW_IN)
	return scroller, nil
}
