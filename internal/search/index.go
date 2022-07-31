// Package search provides an index that allows searching through xkcd comic
// metadata.
package search

import (
	"strconv"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
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

// Load asynchronously fills the comic metadata cache and search index via the
// internet. It may show a loading dialog to the user.
func Load(war WindowAddRemover) error {
	loadingWindow, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return err
	}
	loadingWindow.SetTitle(l("Search Index Update"))
	loadingWindow.SetTypeHint(gdk.WINDOW_TYPE_HINT_DIALOG)
	loadingWindow.SetResizable(false)

	progressBar, err := gtk.ProgressBarNew()
	if err != nil {
		return err
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
		defer func() { done <- struct{}{} }()

		newest, err := cache.NewestComicInfoFromInternet()
		if err != nil {
			return
		}
		for i := 1; i <= newest.Num; i++ {
			n := i
			cache.ComicInfo(n)
			glib.IdleAdd(func() {
				progressBar.SetFraction(float64(n) / float64(newest.Num))
			})
		}
	}()

	// Show cache progress window.
	go func() {
		// Wait before showing the cache progress window. If the cache is
		// already complete, then the caching and indexing operation will be
		// very fast.
		time.Sleep(2 * time.Second)

		select {
		case <-done:
			// Already done, don't bother showing the window.
			glib.IdleAdd(loadingWindow.Destroy)
			return
		default:
			glib.IdleAdd(func() {
				war.AddWindow(loadingWindow)
				loadingWindow.Present()
			})
		}

		// Wait until we are done.
		<-done

		glib.IdleAdd(func() {
			war.RemoveWindow(loadingWindow)
			loadingWindow.Close()
		})
	}()

	return nil
}
