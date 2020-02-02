package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"log"
	"strconv"
)

// AddBookmark adds win's current comic to the user's bookmarks.
func (win *Window) AddBookmark() {
	win.app.bookmarks.Add(win.state.ComicNumber)
}

// RemoveBookmark removes win's current comic from the user's bookmarks.
func (win *Window) RemoveBookmark() {
	win.app.bookmarks.Remove(win.state.ComicNumber)
}

func (win *Window) registerBookmarkObserver() {
	ch := make(chan string)

	win.bookmarkObserverID = win.app.bookmarks.AddObserver(ch)

	go func() {
		for range ch {
			glib.IdleAdd(win.updateBookmarksMenu)
		}
	}()
}

func (win *Window) unregisterBookmarkObserver() {
	win.app.bookmarks.RemoveObserver(win.bookmarkObserverID)
}

func (win *Window) updateBookmarksMenu() {
	win.updateBookmarkButton()
	win.loadBookmarkList()
}

func (win *Window) updateBookmarkButton() {
	if win.app.bookmarks.Contains(win.state.ComicNumber) {
		hasFocus := win.bookmarkActionNew.HasFocus()
		win.actions["bookmark-new"].SetEnabled(false)
		win.bookmarkActionNew.SetVisible(false)

		win.actions["bookmark-remove"].SetEnabled(true)
		win.bookmarkActionRemove.SetVisible(true)
		if hasFocus {
			win.bookmarkActionRemove.GrabFocus()
		}
	} else {
		hasFocus := win.bookmarkActionRemove.HasFocus()
		win.actions["bookmark-remove"].SetEnabled(false)
		win.bookmarkActionRemove.SetVisible(false)

		win.actions["bookmark-new"].SetEnabled(true)
		win.bookmarkActionNew.SetVisible(true)
		if hasFocus {
			win.bookmarkActionNew.GrabFocus()
		}
	}
}

func (win *Window) loadBookmarkList() {
	win.bookmarkList.GetChildren().Foreach(func(child interface{}) {
		win.bookmarkList.Remove(child.(gtk.IWidget))
	})

	if win.app.bookmarks.Empty() {
		win.bookmarkScroller.SetVisible(false)
		win.bookmarkSeparator.SetVisible(false)
		return
	}

	defer win.bookmarkList.ShowAll()
	defer win.bookmarkScroller.SetVisible(true)
	defer win.bookmarkSeparator.SetVisible(true)

	// We are grabbing the newest comic so we can figure out how
	// wide to make the comic number column.
	newest, _ := cache.NewestComicInfo()

	iter := win.app.bookmarks.Iterator()
	for iter.Next() {
		id := iter.Value().(int)
		comic, err := cache.ComicInfo(id)
		if err != nil {
			log.Print("error retrieving comic ", id, ": ", err)
			continue
		}

		item, err := gtk.ButtonNew()
		if err != nil {
			log.Print(err)
			return
		}
		item.Connect("clicked", win.setComicFromBookmark, comic.Num)

		box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
		if err != nil {
			log.Print(err)
			return
		}

		labelID, err := gtk.LabelNew(strconv.Itoa(comic.Num))
		if err != nil {
			log.Print(err)
			return
		}
		labelID.SetXAlign(1)
		labelID.SetWidthChars(len(strconv.Itoa(newest.Num)))
		box.Add(labelID)

		labelTitle, err := gtk.LabelNew(comic.SafeTitle)
		if err != nil {
			log.Print(err)
			return
		}
		labelTitle.SetEllipsize(pango.ELLIPSIZE_END)
		box.Add(labelTitle)

		item.Add(box)
		item.SetRelief(gtk.RELIEF_NONE)
		win.bookmarkList.Add(item)
	}
}

func (win *Window) setComicFromBookmark(_ interface{}, id int) {
	win.SetComic(id)
	win.bookmarks.GetPopover().Hide()
}
