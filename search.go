package main

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
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
}

func (w *Window) UpdateSearch() {
	userQuery, err := w.searchEntry.GetText()
	if err != nil {
		log.Print(err)
	}
	query := query.NewFuzzyQuery(userQuery)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 10
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
	w.searchResults.ShowAll()
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
