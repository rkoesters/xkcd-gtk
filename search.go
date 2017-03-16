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
)

var searchIndex bleve.Index

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
	progressBar, err := gtk.ProgressBarNew()
	if err != nil {
		log.Print(err)
	}
	progressBar.SetText("Building search index...")
	progressBar.SetShowText(true)
	progressBar.Show()
	ca, err := loadingDialog.GetContentArea()
	if err != nil {
		log.Print(err)
	}
	ca.SetMarginTop(12)
	ca.SetMarginBottom(12)
	ca.SetMarginStart(12)
	ca.SetMarginEnd(12)
	ca.Add(progressBar)
	loadingDialog.Present()

	go func() {
		// Make sure all comic metadata is cached and indexed.
		newest, _ := GetNewestComicInfo()
		for i := 1; i <= newest.Num; i++ {
			glib.IdleAdd(func() { progressBar.SetFraction(float64(i) / float64(newest.Num)) })
			GetComicInfo(i)
		}
		glib.IdleAdd(loadingDialog.Close)
	}()
}

func (w *Window) UpdateSearch() {
	userQuery, err := w.searchEntry.GetText()
	if err != nil {
		log.Print(err)
	}
	if userQuery == "" {
		w.clearSearchResults()
		w.loadSearchResults(nil)
		return
	}
	query := query.NewFuzzyQuery(userQuery)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 20
	searchRequest.Fields = []string{"*"}
	result, err := searchIndex.Search(searchRequest)
	if err != nil {
		log.Print(err)
	}
	w.clearSearchResults()
	w.loadSearchResults(result)
}

func (w *Window) clearSearchResults() {
	children := w.searchResults.GetChildren()
	for i := uint(0); i < children.Length(); i++ {
		data := children.NthData(i)
		widget := data.(gtk.IWidget)
		w.searchResults.Remove(widget)
	}
}

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
	for _, sr := range result.Hits {
		item, err := gtk.ButtonNew()
		if err != nil {
			log.Print(err)
			return
		}
		item.Connect("clicked", w.setComicFromSearch, sr.ID)
		label, err := gtk.LabelNew(fmt.Sprintf("%v: %v", sr.ID, sr.Fields["title"]))
		if err != nil {
			log.Print(err)
			return
		}
		label.SetHAlign(gtk.ALIGN_START)
		item.Add(label)
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
func (w *Window) setComicFromSearch(btn *gtk.Button, id string) {
	number, err := strconv.Atoi(id)
	if err != nil {
		log.Print(err)
		return
	}
	w.SetComic(number)
}
