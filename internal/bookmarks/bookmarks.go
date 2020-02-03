// Package bookmarks implements an efficient list for holding a numerically
// sorted list of user bookmarks.
package bookmarks

import (
	"bufio"
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
	"io"
	"os"
	"strconv"
	"sync"
)

// List holds the user's comic bookmarks.
type List struct {
	set *treeset.Set

	observerMutex   sync.RWMutex
	observerCounter int
	observers       map[int]chan string
}

// New returns an initialized List struct.
func New() List {
	return List{
		set: treeset.NewWithIntComparator(),
	}
}

// Add adds the comic number to the bookmarks set.
func (list *List) Add(n int) {
	list.set.Add(n)
	list.notifyObservers("added bookmark " + strconv.Itoa(n))
}

// Remove removes the comic number from the bookmarks set.
func (list *List) Remove(n int) {
	list.set.Remove(n)
	list.notifyObservers("removed bookmark " + strconv.Itoa(n))
}

// Contains indicates whether the comic specified by n is bookmarked.
func (list *List) Contains(n int) bool {
	return list.set.Contains(n)
}

// Empty returns true if there are exactly 0 bookmarks.
func (list *List) Empty() bool {
	return list.set.Empty()
}

// Iterator returns a treeset.Iterator for iterating through the bookmarks.
func (list *List) Iterator() treeset.Iterator {
	return list.set.Iterator()
}

// Read reads bookmarks from r as a newline separated list of comic numbers.
func (list *List) Read(r io.Reader) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return err
		}
		list.Add(n)
	}
	return nil
}

// ReadFile opens the given file and calls Read on the contents.
func (list *List) ReadFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return list.Read(f)
}

// Write writes bookmarks to w as a newline separated list of comic numbers.
func (list *List) Write(w io.Writer) error {
	iter := list.set.Iterator()
	for iter.Next() {
		_, err := fmt.Fprintln(w, iter.Value().(int))
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteFile creates or truncates the given file and calls Write on it.
func (list *List) WriteFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return list.Write(f)
}

// AddObserver adds ch to the list of observers that will be notified when
// changes are made to bookmarks. The returned int can be used to remove the
// added channel from the list of observers using RemoveObserver.
func (list *List) AddObserver(ch chan string) int {
	list.observerMutex.Lock()
	defer list.observerMutex.Unlock()

	if list.observers == nil {
		list.observers = make(map[int]chan string)
	}

	id := list.observerCounter
	list.observerCounter++

	list.observers[id] = ch

	return id
}

// RemoveObserver removes the observer specified by id from the list of
// observers. The channel will be closed after calling this method.
func (list *List) RemoveObserver(id int) {
	list.observerMutex.Lock()
	defer list.observerMutex.Unlock()

	close(list.observers[id])
	delete(list.observers, id)
}

func (list *List) notifyObservers(msg string) {
	list.observerMutex.RLock()
	defer list.observerMutex.RUnlock()

	for _, ch := range list.observers {
		ch <- msg
	}
}
