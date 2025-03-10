// Package search provides an index that allows searching through xkcd comic
// metadata.
package search

import (
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/rkoesters/xkcd"
)

type Index struct {
	index bleve.Index
}

// New initializes and returns a search index. If a search index does not exist
// at the provided path, then New will attempt to create it..
func New(path string) (Index, error) {
	i := Index{}

	var err error
	i.index, err = bleve.Open(path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		i.index, err = bleve.New(path, mapping)
	}
	return i, err
}

// Close closes the search index.
func (i *Index) Close() error {
	return i.index.Close()
}

// Index adds comic to the search index.
func (i *Index) Index(comic *xkcd.Comic) error {
	return i.index.Index(strconv.Itoa(comic.Num), comic)
}

// Search searches the index for the given userQuery.
func (i *Index) Search(userQuery string) (*bleve.SearchResult, error) {
	q := query.NewDisjunctionQuery([]query.Query{
		query.NewQueryStringQuery(userQuery),
		query.NewFuzzyQuery(userQuery),
	})
	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Size = 100
	searchRequest.Fields = []string{"*"}
	return i.index.Search(searchRequest)
}
