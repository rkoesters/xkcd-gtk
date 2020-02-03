package bookmarks_test

import (
	"bytes"
	"github.com/rkoesters/xkcd-gtk/internal/bookmarks"
	"strings"
	"testing"
)

const sortedBookmarkFile = `1
2
3
32
54
432
2345
32456
`

const unsortedBookmarkFile = `1
54
2
3
2345
432
32456
32
`

func TestAddRemove(t *testing.T) {
	bookmarks := bookmarks.New()

	if !bookmarks.Empty() {
		t.Error("New List not empty")
	}

	bookmarks.Add(1)

	if bookmarks.Empty() {
		t.Error("List empty after Add")
	}
	if !bookmarks.Contains(1) {
		t.Error("List does not contain newly added value")
	}

	bookmarks.Remove(1)

	if !bookmarks.Empty() {
		t.Error("List not empty after Remove")
	}
	if bookmarks.Contains(1) {
		t.Error("List Contains removed value")
	}
}

func TestReadWrite(t *testing.T) {
	var buf bytes.Buffer
	bookmarks := bookmarks.New()

	if !bookmarks.Empty() {
		t.Error("New List not empty")
	}

	err := bookmarks.Read(strings.NewReader(sortedBookmarkFile))
	if err != nil {
		t.Fatal(err)
	}

	if bookmarks.Empty() {
		t.Error("List empty after Read")
	}

	err = bookmarks.Write(&buf)
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != sortedBookmarkFile {
		t.Error("Write != Read")
	}
}

func TestReadWriteUnsorted(t *testing.T) {
	var buf bytes.Buffer
	bookmarks := bookmarks.New()

	if !bookmarks.Empty() {
		t.Error("New List not empty")
	}

	err := bookmarks.Read(strings.NewReader(unsortedBookmarkFile))
	if err != nil {
		t.Fatal(err)
	}

	if bookmarks.Empty() {
		t.Error("List empty after Read")
	}

	err = bookmarks.Write(&buf)
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != sortedBookmarkFile {
		t.Error("Write != Read")
	}
}

func TestAddObserver(t *testing.T) {
	notifyCount := 0
	ch := make(chan string)
	done := make(chan struct{})
	go func() {
		for range ch {
			notifyCount++
		}
		done <- struct{}{}
	}()

	bookmarks := bookmarks.New()
	bookmarks.AddObserver(ch)

	if !bookmarks.Empty() {
		t.Error("New List not empty")
	}

	for i := 0; i < 10; i++ {
		bookmarks.Add(i)
	}

	if bookmarks.Empty() {
		t.Error("List empty after 10x Add")
	}

	for i := 0; i < 10; i++ {
		bookmarks.Remove(i)
	}

	if !bookmarks.Empty() {
		t.Error("List not empty after 10x Remove")
	}

	close(ch)
	<-done

	if notifyCount != 20 {
		t.Error("Incorrect notification count")
	}
}

func TestRemoveObserver(t *testing.T) {
	ch := make(chan string)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case _, ok := <-ch:
				if ok {
					t.Error("Received on ch")
				}
			case <-done:
				return
			}
		}
	}()

	bookmarks := bookmarks.New()
	bookmarks.RemoveObserver(bookmarks.AddObserver(ch))

	if !bookmarks.Empty() {
		t.Error("New List not empty")
	}

	for i := 0; i < 10; i++ {
		bookmarks.Add(i)
	}

	if bookmarks.Empty() {
		t.Error("List empty after Add")
	}

	done <- struct{}{}
}
