package widget

import (
	"fmt"
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/search"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type SearchMenu struct {
	*gtk.MenuButton

	popover         *gtk.Popover
	popoverBox      *gtk.Box
	entry           *gtk.SearchEntry
	resultsStack    *gtk.Stack
	resultsNone     *gtk.Label
	resultsScroller *gtk.ScrolledWindow
	resultsList     *ComicListView
}

var _ Widget = &SearchMenu{}

func NewSearchMenu(accels *gtk.AccelGroup, comicSetter func(int)) (*SearchMenu, error) {
	super, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	sm := &SearchMenu{
		MenuButton: super,
	}

	sm.SetTooltipText(l("Search comics"))
	sm.AddAccelerator("activate", accels, gdk.KEY_f, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
	sm.AddAccelerator("activate", accels, gdk.KEY_slash, 0, gtk.ACCEL_VISIBLE)

	sm.popover, err = gtk.PopoverNew(sm)
	if err != nil {
		return nil, err
	}
	sm.SetPopover(sm.popover)

	sm.popoverBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, style.PaddingPopover)
	if err != nil {
		return nil, err
	}
	sm.popoverBox.SetMarginTop(style.PaddingPopover)
	sm.popoverBox.SetMarginBottom(style.PaddingPopover)
	sm.popoverBox.SetMarginStart(style.PaddingPopover)
	sm.popoverBox.SetMarginEnd(style.PaddingPopover)
	sm.entry, err = gtk.SearchEntryNew()
	if err != nil {
		return nil, err
	}
	sm.entry.SetWidthChars(35)
	sm.entry.Connect("search-changed", sm.Search)
	sm.popoverBox.Add(sm.entry)

	sm.resultsStack, err = gtk.StackNew()
	if err != nil {
		return nil, err
	}
	sm.resultsStack.SetHomogeneous(false)
	sm.popoverBox.Add(sm.resultsStack)

	sm.resultsNone, err = gtk.LabelNew(l("No results found"))
	if err != nil {
		return nil, err
	}
	sm.resultsStack.Add(sm.resultsNone)

	sm.resultsScroller, err = NewComicListScroller()
	if err != nil {
		return nil, err
	}
	sm.resultsStack.Add(sm.resultsScroller)

	sm.resultsList, err = NewComicListView(func(n int) {
		comicSetter(n)
		sm.popover.Popdown()
	})
	if err != nil {
		return nil, err
	}
	sm.resultsScroller.Add(sm.resultsList)

	sm.popoverBox.ShowAll()
	sm.popover.Add(sm.popoverBox)

	return sm, sm.loadSearchResults(nil)
}

func (sm *SearchMenu) Dispose() {
	if sm == nil {
		return
	}

	sm.MenuButton = nil

	sm.popover = nil
	sm.popoverBox = nil
	sm.entry = nil
	sm.resultsStack = nil
	sm.resultsNone = nil
	sm.resultsScroller = nil
	sm.resultsList.Dispose()
	sm.resultsList = nil
}

// Search preforms a search with win.searchEntry.GetText() and puts the results
// into win.searchResults.
func (sm *SearchMenu) Search() {
	userQuery, err := sm.entry.GetText()
	if err != nil {
		log.Print("error getting search text: ", err)
	}
	if userQuery == "" {
		err := sm.loadSearchResults(nil)
		if err != nil {
			log.Print("error clearing search results: ", err)
		}
		return
	}
	result, err := search.Search(userQuery)
	if err != nil {
		log.Print("error getting search results: ", err)
	}
	err = sm.loadSearchResults(result)
	if err != nil {
		log.Print("error displaying search results: ", err)
	}
}

// Show the user the given search results.
func (sm *SearchMenu) loadSearchResults(result *bleve.SearchResult) error {
	sm.resultsStack.SetVisible(result != nil)
	if result == nil {
		return nil
	}
	if result.Hits.Len() == 0 {
		sm.resultsStack.SetVisibleChild(sm.resultsNone)
		return nil
	}
	sm.resultsStack.SetVisibleChild(sm.resultsScroller)

	clm, err := NewComicListModel()
	if err != nil {
		return err
	}

	for _, sr := range result.Hits {
		comicNum, err := strconv.Atoi(sr.ID)
		if err != nil {
			return err
		}
		err = clm.AppendComic(comicNum, fmt.Sprint(sr.Fields["safe_title"]))
		if err != nil {
			return err
		}
	}
	sm.resultsList.SetModel(clm)
	return nil
}
