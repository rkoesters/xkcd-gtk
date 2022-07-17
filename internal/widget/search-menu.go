package widget

import (
	"fmt"
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/search"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type SearchMenu struct {
	*gtk.MenuButton

	popover         *gtk.Popover
	popoverBox      *gtk.Box
	entry           *gtk.SearchEntry
	resultsScroller *gtk.ScrolledWindow
	resultsStack    *gtk.Stack
	resultsNone     *gtk.Label
	resultsBox      *gtk.Box

	setComic func(int) // win.SetComic
}

var _ Widget = &SearchMenu{}

func NewSearchMenu(accels *gtk.AccelGroup, comicSetter func(int)) (*SearchMenu, error) {
	super, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	sm := &SearchMenu{
		MenuButton: super,

		setComic: comicSetter,
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
	sm.entry.SetSizeRequest(300, -1)
	sm.entry.Connect("search-changed", sm.Search)
	sm.popoverBox.Add(sm.entry)

	sm.resultsScroller, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	sm.resultsScroller.SetPropagateNaturalHeight(true)
	sm.resultsScroller.SetPropagateNaturalWidth(true)
	sm.resultsScroller.SetMinContentHeight(0)
	sm.resultsScroller.SetMinContentWidth(200)
	sm.resultsScroller.SetMaxContentHeight(350)
	sm.resultsScroller.SetMaxContentWidth(350)
	sm.popoverBox.Add(sm.resultsScroller)

	sm.resultsStack, err = gtk.StackNew()
	if err != nil {
		return nil, err
	}
	sm.resultsStack.SetHomogeneous(false)
	sm.resultsStack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_SLIDE_DOWN)
	sm.resultsScroller.Add(sm.resultsStack)

	sm.resultsNone, err = gtk.LabelNew(l("No results found"))
	if err != nil {
		return nil, err
	}
	sm.resultsStack.Add(sm.resultsNone)

	sm.resultsBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	sm.resultsStack.Add(sm.resultsBox)
	defer func() {
		err := sm.loadSearchResults(nil)
		if err != nil {
			log.Print("error initializing search results: ", err)
		}
	}()

	sm.popoverBox.ShowAll()
	sm.popover.Add(sm.popoverBox)

	return sm, nil
}

func (sm *SearchMenu) Dispose() {
	if sm == nil {
		return
	}

	sm.MenuButton = nil

	sm.popover = nil
	sm.popoverBox = nil
	sm.entry = nil
	sm.resultsScroller = nil
	sm.resultsStack = nil
	sm.resultsNone = nil
	sm.resultsBox = nil
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
	sm.resultsBox.GetChildren().Foreach(func(child interface{}) {
		w, ok := child.(*gtk.Widget)
		if !ok {
			log.Print("error converting child to gtk.Widget")
			return
		}
		sm.resultsBox.Remove(w)
	})

	sm.resultsScroller.SetVisible(result != nil)
	if result == nil {
		return nil
	}
	if result.Hits.Len() == 0 {
		sm.resultsStack.SetVisibleChild(sm.resultsNone)
		return nil
	}
	sm.resultsStack.SetVisibleChild(sm.resultsBox)

	defer sm.resultsBox.ShowAll()

	// We are grabbing the newest comic so we can figure out how wide to
	// make comic ID column.
	newest, err := cache.NewestComicInfoFromCache()
	idWidth := len(strconv.Itoa(newest.Num))
	if err != nil {
		// For the time being, its probably 4 characters.
		idWidth = 4
	}

	for _, sr := range result.Hits {
		id, err := strconv.Atoi(sr.ID)
		if err != nil {
			return err
		}

		clb, err := NewComicListButton(id, fmt.Sprint(sr.Fields["safe_title"]), sm.setComic, idWidth)
		if err != nil {
			return err
		}
		sm.resultsBox.Add(clb)
	}
	return nil
}
