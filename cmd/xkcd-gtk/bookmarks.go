package main

import (
	"github.com/rkoesters/xkcd-gtk/internal/bookmarks"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"log"
	"os"
	"path/filepath"
)

// LoadBookmarks tries to load our bookmarks from disk.
func (app *Application) LoadBookmarks() {
	checkForMisplacedBookmarks()

	app.bookmarks = bookmarks.New()
	app.bookmarks.ReadFile(bookmarksPath())
}

// SaveBookmarks tries to save our bookmarks to disk.
func (app *Application) SaveBookmarks() {
	err := paths.EnsureDataDir()
	if err != nil {
		log.Printf("error saving bookmarks: %v", err)
	}

	err = app.bookmarks.WriteFile(bookmarksPath())
	if err != nil {
		log.Printf("error saving bookmarks: %v", err)
	}
}

func checkForMisplacedBookmarks() {
	misplacedBookmarksList := []string{
		filepath.Join(paths.Builder{}.ConfigDir(), "bookmarks"),
		filepath.Join(paths.Builder{}.DataDir(), "bookmarks"),
		filepath.Join(paths.ConfigDir(), "bookmarks"),
	}

	for _, p := range misplacedBookmarksList {
		_, err := os.Stat(p)
		if !os.IsNotExist(err) {
			log.Printf("WARNING: Potentially misplaced bookmarks file '%v'. Should be '%v'.", p, bookmarksPath())
		}
	}
}

func bookmarksPath() string {
	return filepath.Join(paths.DataDir(), "bookmarks")
}
