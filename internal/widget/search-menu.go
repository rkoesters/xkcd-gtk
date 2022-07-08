package widget

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/search"
	"github.com/rkoesters/xkcd-gtk/internal/style"
	"strconv"
)

type SearchMenu struct {
	menuButton *gtk.MenuButton
	popover    *gtk.Popover
	popoverBox *gtk.Box
	entry      *gtk.SearchEntry
	noResults  *gtk.Label
	results    *gtk.Box
	scroller   *gtk.ScrolledWindow

	setComic func(int) // win.SetComic
}

var _ Widget = &SearchMenu{}

func NewSearchMenu(accels *gtk.AccelGroup, comicSetter func(int)) (*SearchMenu, error) {
	var err error

	sm := &SearchMenu{
		setComic: comicSetter,
	}

	sm.menuButton, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	sm.menuButton.SetTooltipText(l("Search"))
	sm.menuButton.AddAccelerator("activate", accels, gdk.KEY_f, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)

	sm.popover, err = gtk.PopoverNew(sm.menuButton)
	if err != nil {
		return nil, err
	}
	sm.menuButton.SetPopover(sm.popover)

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

	sm.noResults, err = gtk.LabelNew(l("No results found"))
	if err != nil {
		return nil, err
	}
	sm.popoverBox.Add(sm.noResults)

	sm.scroller, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	sm.scroller.SetPropagateNaturalHeight(true)
	sm.scroller.SetPropagateNaturalWidth(true)
	sm.scroller.SetMinContentHeight(0)
	sm.scroller.SetMinContentWidth(200)
	sm.scroller.SetMaxContentHeight(350)
	sm.scroller.SetMaxContentWidth(350)
	sm.popoverBox.Add(sm.scroller)
	sm.results, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	sm.scroller.Add(sm.results)
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

func (sm *SearchMenu) Destroy() {
	sm.menuButton = nil
	sm.popover = nil
	sm.popoverBox = nil
	sm.entry = nil
	sm.noResults = nil
	sm.results = nil
	sm.scroller = nil
}

func (sm *SearchMenu) IWidget() gtk.IWidget {
	return sm.menuButton
}

func (sm *SearchMenu) SetButtonImage(image gtk.IWidget) {
	sm.menuButton.SetImage(image)
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
	sm.results.GetChildren().Foreach(func(child interface{}) {
		sm.results.Remove(child.(gtk.IWidget))
	})

	if result == nil {
		sm.scroller.SetVisible(false)
		sm.noResults.SetVisible(false)
		return nil
	}
	if result.Hits.Len() == 0 {
		sm.scroller.SetVisible(false)
		sm.noResults.SetVisible(true)
		return nil
	}

	defer sm.results.ShowAll()
	defer sm.scroller.SetVisible(true)
	defer sm.noResults.SetVisible(false)

	// We are grabbing the newest comic so we can figure out how wide to
	// make comic ID column.
	newest, err := cache.NewestComicInfoFromCache()
	idWidth := len(strconv.Itoa(newest.Num))
	if err != nil {
		// For the time being, its probably 4 characters.
		idWidth = 4
	}

	for _, sr := range result.Hits {
		item, err := gtk.ModelButtonNew()
		if err != nil {
			return err
		}
		srID := sr.ID
		item.Connect("clicked", func() {
			err := sm.setComicFromSearch(srID)
			if err != nil {
				log.Print("error setting comic from search result: ", err)
			}
		})

		box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, style.PaddingComicListButton)
		if err != nil {
			return err
		}

		labelID, err := gtk.LabelNew(sr.ID)
		if err != nil {
			return err
		}
		labelID.SetXAlign(1)
		// Set character column width using character width of largest
		// comic number.
		labelID.SetWidthChars(idWidth)
		box.Add(labelID)

		labelTitle, err := gtk.LabelNew(fmt.Sprint(sr.Fields["safe_title"]))
		if err != nil {
			return err
		}
		labelTitle.SetEllipsize(pango.ELLIPSIZE_END)
		box.Add(labelTitle)

		child, err := item.GetChild()
		if err != nil {
			return err
		}
		item.Remove(child)
		item.Add(box)
		sm.results.Add(item)
	}
	return nil
}

// setComicFromSearch is a wrapper around win.SetComic to work with search
// result buttons.
func (sm *SearchMenu) setComicFromSearch(id string) error {
	number, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	sm.setComic(number)
	return nil
}
