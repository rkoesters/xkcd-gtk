package widget

import (
	"strconv"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
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

	insertColumn := func(pos int, xalign float64, expand bool, ellipsize pango.EllipsizeMode) (*gtk.TreeViewColumn, error) {
		renderer, err := gtk.CellRendererTextNew()
		if err != nil {
			return nil, err
		}
		renderer.SetAlignment(xalign, 0)
		renderer.SetProperty("xpad", 2)
		renderer.SetProperty("ypad", 6)
		renderer.SetProperty("ellipsize", ellipsize)
		tvc, err := gtk.TreeViewColumnNewWithAttribute(strconv.Itoa(pos), renderer, "text", pos)
		if err != nil {
			return nil, err
		}
		tvc.SetVisible(true)
		tvc.SetExpand(expand)
		tvc.SetSizing(gtk.TREE_VIEW_COLUMN_AUTOSIZE)
		clv.InsertColumn(tvc, pos)
		return tvc, nil
	}

	clv.numberColumn, err = insertColumn(comicListColumnNumber, 1, false, pango.ELLIPSIZE_NONE)
	if err != nil {
		return nil, err
	}

	clv.titleColumn, err = insertColumn(comicListColumnTitle, 0, true, pango.ELLIPSIZE_END)
	if err != nil {
		return nil, err
	}

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

func (clv *ComicListView) rowActivated(tv *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
	itm, err := tv.GetModel()
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
	scroller.SetPropagateNaturalWidth(false) // ComicListView will ellipsize.
	scroller.SetPropagateNaturalHeight(true)
	scroller.SetMaxContentHeight(350)
	scroller.SetShadowType(gtk.SHADOW_IN)
	scroller.SetOverlayScrolling(true)
	return scroller, nil
}
