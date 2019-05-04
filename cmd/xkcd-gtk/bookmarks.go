package main

import (
	"github.com/emirpasic/gods/sets/treeset"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Bookmarks holds the user's comic bookmarks.
type Bookmarks struct {
	set *treeset.Set
}

// Add adds the comic number to the bookmarks set.
func (bookmarks *Bookmarks) Add(n int) {
	bookmarks.set.Add(n)
}

// Remove removes the comic number from the bookmarks set.
func (bookmarks *Bookmarks) Remove(n int) {
	bookmarks.set.Remove(n)
}

// Contains indicates whether the comic specified by n is bookmarked.
func (bookmarks *Bookmarks) Contains(n int) bool {
	return bookmarks.set.Contains(n)
}

// Read reads bookmarks from r.
func (bookmarks *Bookmarks) Read(r io.Reader) error {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return bookmarks.set.FromJSON(bytes)
}

// ReadFile opens the given file and calls Read on the contents.
func (bookmarks *Bookmarks) ReadFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return bookmarks.Read(f)
}

// Write writes bookmarks to w.
func (bookmarks *Bookmarks) Write(w io.Writer) error {
	bytes, err := bookmarks.set.ToJSON()
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

// WriteFile creates or truncates the given file and calls Write on it.
func (bookmarks *Bookmarks) WriteFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return bookmarks.Write(f)
}

// LoadBookmarks tries to load our bookmarks from disk.
func (app *Application) LoadBookmarks() {
	app.bookmarks.set = treeset.NewWithIntComparator()
	app.bookmarks.ReadFile(getBookmarksPath())
}

// SaveBookmarks tries to save our bookmarks to disk.
func (app *Application) SaveBookmarks() {
	err := os.MkdirAll(DataDir(), 0755)
	if err != nil {
		log.Printf("error saving bookmarks: %v", err)
	}

	err = app.bookmarks.WriteFile(getBookmarksPath())
	if err != nil {
		log.Printf("error saving bookmarks: %v", err)
	}
}

func getBookmarksPath() string {
	return filepath.Join(DataDir(), "bookmarks")
}
