package main

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"log"
	"path/filepath"
	"strconv"
	"time"
)

var searchIndex bleve.Index

func initSearchIndex() error {
	var err error

	searchIndex, err = bleve.Open(getSearchIndexPath())
	if err == bleve.ErrorIndexPathDoesNotExist {
		// searchIndex doesn't exist yet, lets make it.
		mapping := bleve.NewIndexMapping()
		searchIndex, err = bleve.New(getSearchIndexPath(), mapping)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func closeSearchIndex() error {
	return searchIndex.Close()
}

// LoadSearchIndex makes sure that every xkcd comic metadata is cached and
// indexed in the search index.
func (app *Application) LoadSearchIndex() {
	loadingWindow, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Print(err)
	}
	loadingWindow.SetTitle(l("Search Index Update"))
	loadingWindow.SetTypeHint(gdk.WINDOW_TYPE_HINT_DIALOG)
	loadingWindow.SetResizable(false)

	progressBar, err := gtk.ProgressBarNew()
	if err != nil {
		log.Print(err)
	}
	progressBar.SetText(l("Updating comic search index..."))
	progressBar.SetShowText(true)
	progressBar.SetMarginTop(24)
	progressBar.SetMarginBottom(24)
	progressBar.SetMarginStart(24)
	progressBar.SetMarginEnd(24)
	progressBar.SetSizeRequest(300, -1)
	progressBar.SetFraction(0)
	progressBar.Show()
	loadingWindow.Add(progressBar)

	done := make(chan struct{})

	// Make sure all comic metadata is cached and indexed.
	go func() {
		newest, _ := GetNewestComicInfo()
		for i := 1; i <= newest.Num; i++ {
			n := i
			GetComicInfo(n)
			glib.IdleAdd(func() {
				progressBar.SetFraction(float64(n) / float64(newest.Num))
			})
		}
		done <- struct{}{}
	}()

	// Show cache progress window.
	go func() {
		// Wait before showing the cache progress window. If the cache
		// is already complete, then the caching and indexing operation
		// will be very fast.
		time.Sleep(time.Second)

		select {
		case <-done:
			// Already done, don't bother showing the window.
			glib.IdleAdd(loadingWindow.Destroy)
			return
		default:
			glib.IdleAdd(func() {
				app.application.AddWindow(loadingWindow)
				loadingWindow.Present()
			})
		}

		// Wait until we are done.
		<-done

		glib.IdleAdd(func() {
			app.application.RemoveWindow(loadingWindow)
			loadingWindow.Close()
		})
	}()
}

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
	query := query.NewQueryStringQuery(userQuery)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 50
	searchRequest.Fields = []string{"*"}
	result, err := searchIndex.Search(searchRequest)
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
	// make comic Id column.
	newest, _ := GetNewestComicInfo()

	for _, sr := range result.Hits {
		item, err := gtk.ButtonNew()
		if err != nil {
			log.Print(err)
			return
		}
		item.Connect("clicked", win.setComicFromSearch, sr.ID)

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
func (win *Window) setComicFromSearch(_ interface{}, id string) {
	number, err := strconv.Atoi(id)
	if err != nil {
		log.Print("error setting comic from search result: ", err)
		return
	}
	win.SetComic(number)
	win.search.GetPopover().Hide()
}

func getSearchIndexPath() string {
	return filepath.Join(paths.CacheDir(), "search")
}
