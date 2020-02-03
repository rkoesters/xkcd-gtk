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
	app.bookmarks = bookmarks.New()
	app.bookmarks.ReadFile(bookmarksPath())
}

// SaveBookmarks tries to save our bookmarks to disk.
func (app *Application) SaveBookmarks() {
	err := os.MkdirAll(paths.DataDir(), 0755)
	if err != nil {
		log.Printf("error saving bookmarks: %v", err)
	}

	err = app.bookmarks.WriteFile(bookmarksPath())
	if err != nil {
		log.Printf("error saving bookmarks: %v", err)
	}
}

func bookmarksPath() string {
	return filepath.Join(paths.DataDir(), "bookmarks")
}
