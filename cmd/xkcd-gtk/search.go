package main

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/search"
	"log"
	"strconv"
)

// Search preforms a search with win.searchEntry.GetText() and puts the results
// into win.searchResults.
func (win *Window) Search() {
	userQuery, err := win.searchEntry.GetText()
	if err != nil {
		log.Print("error getting search text: ", err)
	}
	if userQuery == "" {
		win.loadSearchResults(nil)
		return
	}
	result, err := search.Search(userQuery)
	if err != nil {
		log.Print("error getting search results: ", err)
	}
	win.loadSearchResults(result)
}

// Show the user the given search results.
func (win *Window) loadSearchResults(result *bleve.SearchResult) {
	win.searchResults.GetChildren().Foreach(func(child interface{}) {
		win.searchResults.Remove(child.(gtk.IWidget))
	})

	if result == nil {
		win.searchScroller.SetVisible(false)
		win.searchNoResults.SetVisible(false)
		return
	}
	if result.Hits.Len() == 0 {
		win.searchScroller.SetVisible(false)
		win.searchNoResults.SetVisible(true)
		return
	}

	defer win.searchResults.ShowAll()
	defer win.searchScroller.SetVisible(true)
	defer win.searchNoResults.SetVisible(false)

	// We are grabbing the newest comic so we can figure out how wide to
	// make comic ID column.
	newest, _ := cache.NewestComicInfo()

	for _, sr := range result.Hits {
		item, err := gtk.ButtonNew()
		if err != nil {
			log.Print(err)
			return
		}
		srID := sr.ID
		item.Connect("clicked", func() { win.setComicFromSearch(srID) })

		box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
		if err != nil {
			log.Print(err)
			return
		}

		labelID, err := gtk.LabelNew(sr.ID)
		if err != nil {
			log.Print(err)
			return
		}
		labelID.SetXAlign(1)
		// Set character column width using character width of largest
		// comic number.
		labelID.SetWidthChars(len(strconv.Itoa(newest.Num)))
		box.Add(labelID)

		labelTitle, err := gtk.LabelNew(fmt.Sprint(sr.Fields["safe_title"]))
		if err != nil {
			log.Print(err)
			return
		}
		labelTitle.SetEllipsize(pango.ELLIPSIZE_END)
		box.Add(labelTitle)

		item.Add(box)
		item.SetRelief(gtk.RELIEF_NONE)
		win.searchResults.Add(item)
	}
}

// setComicFromSearch is a wrapper around win.SetComic to work with search
// result buttons.
func (win *Window) setComicFromSearch(id string) {
	number, err := strconv.Atoi(id)
	if err != nil {
		log.Print("error setting comic from search result: ", err)
		return
	}
	win.SetComic(number)
	win.search.GetPopover().Hide()
}
