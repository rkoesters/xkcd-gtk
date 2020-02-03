package main

import (
	"bufio"
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Bookmarks holds the user's comic bookmarks.
type Bookmarks struct {
	set *treeset.Set

	observerMutex   sync.Mutex
	observerCounter int
	observers       map[int]chan string
}

// Add adds the comic number to the bookmarks set.
func (bookmarks *Bookmarks) Add(n int) {
	bookmarks.set.Add(n)
	bookmarks.notifyObservers("added bookmark " + strconv.Itoa(n))
}

// Remove removes the comic number from the bookmarks set.
func (bookmarks *Bookmarks) Remove(n int) {
	bookmarks.set.Remove(n)
	bookmarks.notifyObservers("removed bookmark " + strconv.Itoa(n))
}

// Contains indicates whether the comic specified by n is bookmarked.
func (bookmarks *Bookmarks) Contains(n int) bool {
	return bookmarks.set.Contains(n)
}

// Empty returns true if there are exactly 0 bookmarks.
func (bookmarks *Bookmarks) Empty() bool {
	return bookmarks.set.Empty()
}

// Iterator returns a treeset.Iterator for iterating through the bookmarks.
func (bookmarks *Bookmarks) Iterator() treeset.Iterator {
	return bookmarks.set.Iterator()
}

// Read reads bookmarks from r as a newline separated list of comic numbers.
func (bookmarks *Bookmarks) Read(r io.Reader) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return err
		}
		bookmarks.Add(n)
	}
	return nil
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

// Write writes bookmarks to w as a newline separated list of comic numbers.
func (bookmarks *Bookmarks) Write(w io.Writer) error {
	iter := bookmarks.set.Iterator()
	for iter.Next() {
		_, err := fmt.Fprintln(w, iter.Value().(int))
		if err != nil {
			return err
		}
	}
	return nil
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

// AddObserver adds ch to the list of observers that will be notified when
// changes are made to bookmarks. The returned int can be used to remove the
// added channel from the list of observers using RemoveObserver.
func (bookmarks *Bookmarks) AddObserver(ch chan string) int {
	bookmarks.observerMutex.Lock()
	defer bookmarks.observerMutex.Unlock()

	if bookmarks.observers == nil {
		bookmarks.observers = make(map[int]chan string)
	}

	id := bookmarks.observerCounter
	bookmarks.observerCounter++

	bookmarks.observers[id] = ch

	return id
}

// RemoveObserver removes the observer specified by id from the list of
// observers. The channel will be closed after calling this method.
func (bookmarks *Bookmarks) RemoveObserver(id int) {
	bookmarks.observerMutex.Lock()
	defer bookmarks.observerMutex.Unlock()

	close(bookmarks.observers[id])
	delete(bookmarks.observers, id)
}

func (bookmarks *Bookmarks) notifyObservers(msg string) {
	bookmarks.observerMutex.Lock()
	defer bookmarks.observerMutex.Unlock()

	for _, ch := range bookmarks.observers {
		ch <- msg
	}
}

// LoadBookmarks tries to load our bookmarks from disk.
func (app *Application) LoadBookmarks() {
	app.bookmarks.set = treeset.NewWithIntComparator()
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
