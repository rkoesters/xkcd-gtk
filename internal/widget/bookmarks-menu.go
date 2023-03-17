package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/bookmarks"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type BookmarksMenu struct {
	*gtk.ButtonBox

	bookmarkButton *gtk.Button

	popoverButton *gtk.MenuButton
	popover       *gtk.Popover
	popoverBox    *gtk.Box
	scroller      *gtk.ScrolledWindow
	list          *ComicListView

	bookmarks *bookmarks.List               // ptr to app.bookmarks
	actions   map[string]*glib.SimpleAction // ptr to win.actions

	updateButtonIcons func()
}

var _ Widget = &BookmarksMenu{}

func NewBookmarksMenu(b *bookmarks.List, actions map[string]*glib.SimpleAction, accels *gtk.AccelGroup, comicSetter func(int), updateButtonIcons func()) (*BookmarksMenu, error) {
	super, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	bm := &BookmarksMenu{
		ButtonBox: super,

		bookmarks: b,
		actions:   actions,

		updateButtonIcons: updateButtonIcons,
	}
	bm.SetLayout(gtk.BUTTONBOX_EXPAND)
	bm.SetHomogeneous(false)

	bm.bookmarkButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	bm.bookmarkButton.AddAccelerator("activate", accels, gdk.KEY_d, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
	bm.Add(bm.bookmarkButton)

	bm.popoverButton, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	bm.popoverButton.SetTooltipText(l("Bookmarks"))
	bm.popoverButton.AddAccelerator("activate", accels, gdk.KEY_b, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
	bm.popoverButton.AddAccelerator("activate", accels, gdk.KEY_d, gdk.CONTROL_MASK|gdk.MOD1_MASK, gtk.ACCEL_VISIBLE)
	bm.Add(bm.popoverButton)

	bm.popover, err = gtk.PopoverNew(bm.popoverButton)
	if err != nil {
		return nil, err
	}
	bm.popoverButton.SetPopover(bm.popover)

	bm.popoverBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, style.PaddingPopover)
	if err != nil {
		return nil, err
	}
	bm.popoverBox.SetMarginTop(style.PaddingPopover)
	bm.popoverBox.SetMarginBottom(style.PaddingPopover)
	bm.popoverBox.SetMarginStart(style.PaddingPopover)
	bm.popoverBox.SetMarginEnd(style.PaddingPopover)

	bm.scroller, err = NewComicListScroller()
	if err != nil {
		return nil, err
	}
	bm.popoverBox.Add(bm.scroller)

	bm.list, err = NewComicListView(func(n int) {
		comicSetter(n)
		bm.popover.Popdown()
	})
	if err != nil {
		return nil, err
	}
	bm.list.SetSizeRequest(280, -1)
	bm.scroller.Add(bm.list)

	defer func() {
		err := bm.loadBookmarkList()
		if err != nil {
			log.Print("error calling loadBookmarkList(): ", err)
		}
	}()

	bm.popoverBox.ShowAll()
	bm.popover.Add(bm.popoverBox)

	sc, err := bm.GetStyleContext()
	if err != nil {
		return nil, err
	}
	sc.AddClass(style.ClassLinked)
	bm.SetSpacing(0)

	sc, err = bm.popoverButton.GetStyleContext()
	if err != nil {
		return nil, err
	}
	sc.AddClass(style.ClassSlimButton)

	return bm, nil
}

func (bm *BookmarksMenu) Dispose() {
	if bm == nil {
		return
	}

	bm.ButtonBox = nil

	bm.bookmarkButton = nil

	bm.popoverButton = nil
	bm.popover = nil
	bm.popoverBox = nil
	bm.scroller = nil
	bm.list.Dispose()
	bm.list = nil

	bm.bookmarks = nil
	bm.actions = nil
}

func (bm *BookmarksMenu) Update(comicNumber int) {
	err := bm.loadBookmarkList()
	if err != nil {
		log.Print("error calling loadBookmarkList(): ", err)
	}

	bookmarked := bm.bookmarks.Contains(comicNumber)

	bm.actions["bookmark-new"].SetEnabled(!bookmarked)
	bm.actions["bookmark-remove"].SetEnabled(bookmarked)

	if bookmarked {
		bm.bookmarkButton.SetActionName("win.bookmark-remove")
		bm.bookmarkButton.SetTooltipText(l("Remove from bookmarks"))
	} else {
		bm.bookmarkButton.SetActionName("win.bookmark-new")
		bm.bookmarkButton.SetTooltipText(l("Add to bookmarks"))
	}

	glib.IdleAdd(bm.updateButtonIcons)
}

func (bm *BookmarksMenu) loadBookmarkList() error {
	empty := bm.bookmarks.Empty()
	bm.popoverButton.SetSensitive(!empty)
	if empty {
		return nil
	}

	clm, err := NewComicListModel()
	if err != nil {
		return err
	}

	iter := bm.bookmarks.Iterator()
	for iter.Next() {
		comicNumber := iter.Value().(int)
		comic, err := cache.ComicInfo(comicNumber)
		if err != nil {
			return err
		}
		err = clm.AppendComic(comicNumber, comic.SafeTitle)
		if err != nil {
			return err
		}
	}
	bm.list.SetModel(clm)
	return nil
}

func (bm *BookmarksMenu) SetLinkedButtons(linked bool) error {
	sc, err := bm.GetStyleContext()
	if err != nil {
		return err
	}

	if linked {
		sc.AddClass(style.ClassLinked)
	} else {
		sc.RemoveClass(style.ClassLinked)
	}

	return nil
}
