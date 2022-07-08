package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/bookmarks"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/style"
	"strconv"
)

type BookmarksMenu struct {
	bookmarks  *bookmarks.List // ptr to app.bookmarks
	observerID int

	windowState *WindowState                  // ptr to win.state
	actions     map[string]*glib.SimpleAction // ptr to win.actions

	menuButton   *gtk.MenuButton
	popover      *gtk.Popover
	popoverBox   *gtk.Box
	addButton    *gtk.Button
	removeButton *gtk.Button
	separator    *gtk.Separator
	scroller     *gtk.ScrolledWindow
	list         *gtk.Box

	setComic func(int) // win.SetComic
}

var _ Widget = &BookmarksMenu{}

func NewBookmarksMenu(b *bookmarks.List, win *gtk.ApplicationWindow, ws *WindowState, actions map[string]*glib.SimpleAction, accels *gtk.AccelGroup, comicSetter func(int)) (*BookmarksMenu, error) {
	var err error

	bm := &BookmarksMenu{
		bookmarks:   b,
		windowState: ws,
		actions:     actions,
		setComic:    comicSetter,
	}

	// Create the bookmark menu
	bm.menuButton, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	bm.menuButton.SetTooltipText(l("Bookmarks"))
	bm.menuButton.AddAccelerator("activate", accels, gdk.KEY_b, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)

	bm.popover, err = gtk.PopoverNew(bm.menuButton)
	if err != nil {
		return nil, err
	}
	bm.menuButton.SetPopover(bm.popover)

	bm.popoverBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
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
	bookmarkNewImage, err := gtk.ImageNewFromIconName("bookmark-new-symbolic", gtk.ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}
	bm.addButton.SetImage(bookmarkNewImage)
	bm.addButton.SetAlwaysShowImage(true)
	bm.popoverBox.Add(bm.addButton)

	bm.removeButton, err = gtk.ButtonNewWithLabel(l("Remove from bookmarks"))
	if err != nil {
		return nil, err
	}
	bm.removeButton.SetActionName("win.bookmark-remove")
	bookmarkRemoveImage, err := gtk.ImageNewFromIconName("edit-delete-symbolic", gtk.ICON_SIZE_BUTTON)
	if err != nil {
		return nil, err
	}
	bm.removeButton.SetImage(bookmarkRemoveImage)
	bm.removeButton.SetAlwaysShowImage(true)
	bm.popoverBox.Add(bm.removeButton)

	bm.separator, err = gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	bm.popoverBox.PackStart(bm.separator, false, false, style.PaddingPopover)

	bm.scroller, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	bm.scroller.SetPropagateNaturalHeight(true)
	bm.scroller.SetPropagateNaturalWidth(true)
	bm.scroller.SetMinContentHeight(0)
	bm.scroller.SetMinContentWidth(200)
	bm.scroller.SetMaxContentHeight(350)
	bm.scroller.SetMaxContentWidth(350)
	bm.popoverBox.Add(bm.scroller)
	bm.list, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	bm.scroller.Add(bm.list)
	bm.registerBookmarkObserver()
	win.Connect("delete-event", func() {
		bm.unregisterBookmarkObserver()
	})
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

func (bm *BookmarksMenu) Destroy() {
	bm.bookmarks = nil

	bm.windowState = nil
	bm.actions = nil

	bm.menuButton = nil
	bm.popover = nil
	bm.popoverBox = nil
	bm.addButton = nil
	bm.removeButton = nil
	bm.separator = nil
	bm.scroller = nil
	bm.list = nil
}

func (bm *BookmarksMenu) IWidget() gtk.IWidget {
	return bm.menuButton
}

func (bm *BookmarksMenu) SetButtonImage(image gtk.IWidget) {
	bm.menuButton.SetImage(image)
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
	if bm.bookmarks.Contains(bm.windowState.ComicNumber) {
		hasFocus := bm.addButton.HasFocus()
		bm.actions["bookmark-new"].SetEnabled(false)
		bm.addButton.SetVisible(false)

		bm.actions["bookmark-remove"].SetEnabled(true)
		bm.removeButton.SetVisible(true)
		if hasFocus {
			bm.removeButton.GrabFocus()
		}
	} else {
		hasFocus := bm.removeButton.HasFocus()
		bm.actions["bookmark-remove"].SetEnabled(false)
		bm.removeButton.SetVisible(false)

		bm.actions["bookmark-new"].SetEnabled(true)
		bm.addButton.SetVisible(true)
		if hasFocus {
			bm.addButton.GrabFocus()
		}
	}
}

func (bm *BookmarksMenu) loadBookmarkList() error {
	bm.list.GetChildren().Foreach(func(child interface{}) {
		bm.list.Remove(child.(gtk.IWidget))
	})

	if bm.bookmarks.Empty() {
		bm.scroller.SetVisible(false)
		bm.separator.SetVisible(false)
		return nil
	}

	defer bm.list.ShowAll()
	defer bm.scroller.SetVisible(true)
	defer bm.separator.SetVisible(true)

	// We are grabbing the newest comic so we can figure out how
	// wide to make the comic number column.
	newest, err := cache.NewestComicInfoFromCache()
	idWidth := len(strconv.Itoa(newest.Num))
	if err != nil {
		// For the time being, its probably 4 characters.
		idWidth = 4
	}

	iter := bm.bookmarks.Iterator()
	for iter.Next() {
		id := iter.Value().(int)
		comic, err := cache.ComicInfo(id)
		if err != nil {
			return err
		}
		clb, err := NewComicListButton(comic.Num, comic.SafeTitle, bm.setComic, idWidth)
		if err != nil {
			return err
		}
		bm.list.Add(clb.IWidget())
	}
	return nil
}
