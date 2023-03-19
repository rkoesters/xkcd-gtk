// Package search provides an index that allows searching through xkcd comic
// metadata.
package search

import (
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
)

var index bleve.Index

// Init initializes the search index.
func Init() (err error) {
	paths.CheckForMisplacedSearchIndex()

	log.Debug("opening search index: ", paths.SearchIndex())
	index, err = bleve.Open(paths.SearchIndex())
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Debug("search index not found, creating new search index")
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(paths.SearchIndex(), mapping)
	}
	return
}

// Close closes the search index.
func Close() error {
	return index.Close()
}

// Index adds comic to the search index.
func Index(comic *xkcd.Comic) error {
	log.Debug("indexing: ", comic)
	return index.Index(strconv.Itoa(comic.Num), comic)
}

// Search searches the index for the given userQuery.
func Search(userQuery string) (*bleve.SearchResult, error) {
	q := query.NewDisjunctionQuery([]query.Query{
		query.NewQueryStringQuery(userQuery),
		query.NewFuzzyQuery(userQuery),
	})
	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Size = 100
	searchRequest.Fields = []string{"*"}
	return index.Search(searchRequest)
}

type WindowAddRemover interface {
	AddWindow(gtk.IWindow)
	RemoveWindow(gtk.IWindow)
}
