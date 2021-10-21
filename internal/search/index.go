// Package search provides an index that allows searching through xkcd comic
// metadata.
package search

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var index bleve.Index

// Init initializes the search index.
func Init() (err error) {
	checkForMisplacedSearchIndex()

	log.Debug("opening search index: ", searchIndexPath())
	index, err = bleve.Open(searchIndexPath())
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Debug("search index not found, creating new search index")
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(searchIndexPath(), mapping)
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
	q := query.NewQueryStringQuery(userQuery)
	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Size = 50
	searchRequest.Fields = []string{"*"}
	return index.Search(searchRequest)
}

// Load asynchronously fills the comic metadata cache and search index via the
// internet. It may show a loading dialog to the user.
func Load(app *gtk.Application) {
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
				app.AddWindow(loadingWindow)
				loadingWindow.Present()
			})
		}

		// Wait until we are done.
		<-done

		glib.IdleAdd(func() {
			app.RemoveWindow(loadingWindow)
			loadingWindow.Close()
		})
	}()
}

func checkForMisplacedSearchIndex() {
	misplacedSearchIndex := filepath.Join(paths.Builder{}.CacheDir(), "search")

	_, err := os.Stat(misplacedSearchIndex)
	if !os.IsNotExist(err) {
		log.Printf("WARNING: Potentially misplaced search index '%v'. Should be '%v'.", misplacedSearchIndex, searchIndexPath())
	}
}

func searchIndexPath() string {
	return filepath.Join(paths.CacheDir(), "search")
}
