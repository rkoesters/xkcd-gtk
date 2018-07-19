package main

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"path/filepath"
	"strconv"
	"time"
)

var searchIndex bleve.Index

// LoadSearchIndex makes sure that every xkcd comic metadata is cached
// and indexed in the search index.
func (a *Application) LoadSearchIndex() {
	var err error
	searchIndexPath := filepath.Join(CacheDir(), "search")
	searchIndex, err = bleve.Open(searchIndexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		// searchIndex doesn't exist yet, lets make it.
		mapping := bleve.NewIndexMapping()
		searchIndex, err = bleve.New(searchIndexPath, mapping)
		if err != nil {
			log.Print(err)
		}
	} else if err != nil {
		log.Print(err)
	}

	loadingDialog, err := gtk.DialogNew()
	if err != nil {
		log.Print(err)
	}
	loadingDialog.SetTitle("Search Index Update")
	loadingDialog.SetResizable(false)
	progressBar, err := gtk.ProgressBarNew()
	if err != nil {
		log.Print(err)
	}
	progressBar.SetText("Updating comic search index...")
	progressBar.SetShowText(true)
	progressBar.SetFraction(0)
	progressBar.Show()
	ca, err := loadingDialog.GetContentArea()
	if err != nil {
		log.Print(err)
	}
	// A gtk.Dialog content area has some children by default, we want
	// to remove those children so the only child is progressBar.
	ca.GetChildren().Foreach(func(child interface{}) {
		ca.Remove(child.(gtk.IWidget))
	})
	ca.SetMarginTop(24)
	ca.SetMarginBottom(24)
	ca.SetMarginStart(24)
	ca.SetMarginEnd(24)
	ca.Add(progressBar)

	done := false

	// Lets only open the dialog if our loading will be longer.
	go func() {
		time.Sleep(time.Second)
		if !done {
			glib.IdleAdd(loadingDialog.Present)
		}
	}()

	// Make sure all comic metadata is cached and indexed.
	go func() {
		newest, _ := GetNewestComicInfo()
		for i := 1; i <= newest.Num; i++ {
			GetComicInfo(i)
			glib.IdleAdd(func() { progressBar.SetFraction(float64(i) / float64(newest.Num)) })
		}
		done = true
		glib.IdleAdd(loadingDialog.Close)
	}()
}

// Search preforms a search with w.searchEntry.GetText() and puts the
// results into w.searchResults.
func (w *Window) Search() {
	userQuery, err := w.searchEntry.GetText()
	if err != nil {
		log.Print(err)
	}
	if userQuery == "" {
		w.clearSearchResults()
		w.loadSearchResults(nil)
		return
	}
	query := query.NewQueryStringQuery(userQuery)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 50
	searchRequest.Fields = []string{"*"}
	result, err := searchIndex.Search(searchRequest)
	if err != nil {
		log.Print(err)
	}
	w.clearSearchResults()
	w.loadSearchResults(result)
}

// Remove all widgets from the search results area.
func (w *Window) clearSearchResults() {
	w.searchResults.GetChildren().Foreach(func(child interface{}) {
		w.searchResults.Remove(child.(gtk.IWidget))
	})
}

// Show the user the given search results.
func (w *Window) loadSearchResults(result *bleve.SearchResult) {
	defer w.searchResults.ShowAll()
	if result == nil {
		// If there are no results to display, show a friendly message.
		label, err := gtk.LabelNew("Whatcha lookin' for?")
		if err != nil {
			log.Print(err)
			return
		}
		label.SetVExpand(true)
		w.searchResults.Add(label)
		return
	}
	// We are grabbing the newest comic so we can figure out how wide to
	// make comic Id column.
	newest, _ := GetNewestComicInfo()
	for _, sr := range result.Hits {
		item, err := gtk.ButtonNew()
		if err != nil {
			log.Print(err)
			return
		}
		item.Connect("clicked", w.setComicFromSearch, sr.ID)
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
		labelID.SetWidthChars(len(fmt.Sprint(newest.Num)))
		box.Add(labelID)
		labelTitle, err := gtk.LabelNew(fmt.Sprint(sr.Fields["safe_title"]))
		if err != nil {
			log.Print(err)
			return
		}
		box.Add(labelTitle)
		item.Add(box)
		item.SetRelief(gtk.RELIEF_NONE)
		w.searchResults.Add(item)
	}
	if result.Hits.Len() == 0 {
		label, err := gtk.LabelNew("0 search results")
		if err != nil {
			log.Print(err)
			return
		}
		label.SetVExpand(true)
		w.searchResults.Add(label)
	}
}

// setComicFromSearch is a wrapper around w.SetComic to work with search
// result buttons.
func (w *Window) setComicFromSearch(_ interface{}, id string) {
	number, err := strconv.Atoi(id)
	if err != nil {
		log.Print(err)
		return
	}
	w.SetComic(number)
}