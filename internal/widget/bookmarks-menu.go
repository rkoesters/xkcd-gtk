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
	*gtk.MenuButton

	popover      *gtk.Popover
	popoverBox   *gtk.Box
	addButton    *gtk.Button
	removeButton *gtk.Button
	scroller     *gtk.ScrolledWindow
	list         *ComicListView

	bookmarks  *bookmarks.List // ptr to app.bookmarks
	observerID int

	windowState *WindowState                  // ptr to win.state
	actions     map[string]*glib.SimpleAction // ptr to win.actions
}

var _ Widget = &BookmarksMenu{}

func NewBookmarksMenu(b *bookmarks.List, win *ApplicationWindow, ws *WindowState, actions map[string]*glib.SimpleAction, accels *gtk.AccelGroup, comicSetter func(int)) (*BookmarksMenu, error) {
	const (
		bmIconSize = gtk.ICON_SIZE_MENU
		btnWidth   = 280
		btnHeight  = -1
	)

	super, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	bm := &BookmarksMenu{
		MenuButton: super,

		bookmarks:   b,
		windowState: ws,
		actions:     actions,
	}

	bm.SetTooltipText(l("Bookmarks"))
	bm.AddAccelerator("activate", accels, gdk.KEY_b, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)

	bm.popover, err = gtk.PopoverNew(bm)
	if err != nil {
		return nil, err
	}
	bm.SetPopover(bm.popover)

	bm.popoverBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, style.PaddingPopover)
	if err != nil {
		return nil, err
	}
	bm.popoverBox.SetMarginTop(style.PaddingPopover)
	bm.popoverBox.SetMarginBottom(style.PaddingPopover)
	bm.popoverBox.SetMarginStart(style.PaddingPopover)
	bm.popoverBox.SetMarginEnd(style.PaddingPopover)

	bm.addButton, err = gtk.ButtonNewWithLabel(l("Add to bookmarks"))
	if err != nil {
		return nil, err
	}
	bm.addButton.SetActionName("win.bookmark-new")
	bookmarkNewImage, err := gtk.ImageNewFromIconName("bookmark-new-symbolic", bmIconSize)
	if err != nil {
		return nil, err
	}
	bm.addButton.SetImage(bookmarkNewImage)
	bm.addButton.SetAlwaysShowImage(true)
	bm.addButton.SetSizeRequest(btnWidth, btnHeight)
	bm.popoverBox.Add(bm.addButton)

	bm.removeButton, err = gtk.ButtonNewWithLabel(l("Remove from bookmarks"))
	if err != nil {
		return nil, err
	}
	bm.removeButton.SetActionName("win.bookmark-remove")
	bookmarkRemoveImage, err := gtk.ImageNewFromIconName("edit-delete-symbolic", bmIconSize)
	if err != nil {
		return nil, err
	}
	bm.removeButton.SetImage(bookmarkRemoveImage)
	bm.removeButton.SetAlwaysShowImage(true)
	bm.removeButton.SetSizeRequest(btnWidth, btnHeight)
	bm.popoverBox.Add(bm.removeButton)

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
	bm.scroller.Add(bm.list)

	bm.registerBookmarkObserver()
	win.Connect("delete-event", bm.unregisterBookmarkObserver)
	defer func() {
		err := bm.loadBookmarkList()
		if err != nil {
			log.Print("error calling loadBookmarkList(): ", err)
		}
	}()

	bm.popoverBox.ShowAll()
	bm.popover.Add(bm.popoverBox)

	return bm, nil
}

func (bm *BookmarksMenu) Dispose() {
	if bm == nil {
		return
	}

	bm.MenuButton = nil

	bm.popover = nil
	bm.popoverBox = nil
	bm.addButton = nil
	bm.removeButton = nil
	bm.scroller = nil
	bm.list.Dispose()
	bm.list = nil

	bm.bookmarks = nil

	bm.windowState = nil
	bm.actions = nil
}

// AddBookmark adds win's current comic to the user's bookmarks.
func (bm *BookmarksMenu) AddBookmark() {
	bm.bookmarks.Add(bm.windowState.ComicNumber)
}

// RemoveBookmark removes win's current comic from the user's bookmarks.
func (bm *BookmarksMenu) RemoveBookmark() {
	bm.bookmarks.Remove(bm.windowState.ComicNumber)
}

func (bm *BookmarksMenu) registerBookmarkObserver() {
	ch := make(chan string)

	bm.observerID = bm.bookmarks.AddObserver(ch)

	go func() {
		for range ch {
			glib.IdleAdd(bm.UpdateBookmarksMenu)
		}
	}()
}

func (bm *BookmarksMenu) unregisterBookmarkObserver() {
	bm.bookmarks.RemoveObserver(bm.observerID)
}

func (bm *BookmarksMenu) UpdateBookmarksMenu() {
	bm.UpdateBookmarkButton()
	err := bm.loadBookmarkList()
	if err != nil {
		log.Print("error calling loadBookmarkList(): ", err)
	}
}

func (bm *BookmarksMenu) UpdateBookmarkButton() {
	currentIsBookmarked := bm.bookmarks.Contains(bm.windowState.ComicNumber)

	var focused bool
	if currentIsBookmarked {
		focused = bm.addButton.IsFocus()
	} else {
		focused = bm.removeButton.IsFocus()
	}

	bm.addButton.SetVisible(!currentIsBookmarked)
	bm.removeButton.SetVisible(currentIsBookmarked)
	bm.actions["bookmark-new"].SetEnabled(!currentIsBookmarked)
	bm.actions["bookmark-remove"].SetEnabled(currentIsBookmarked)

	if focused {
		if currentIsBookmarked {
			bm.removeButton.GrabFocus()
		} else {
			bm.addButton.GrabFocus()
		}
	}
}

func (bm *BookmarksMenu) loadBookmarkList() error {
	empty := bm.bookmarks.Empty()
	bm.scroller.SetVisible(!empty)
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
