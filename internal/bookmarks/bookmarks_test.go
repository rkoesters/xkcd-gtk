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

func TestReadWrite(t *testing.T) {
	var buf bytes.Buffer
	bookmarks := bookmarks.New()

	err := bookmarks.Read(strings.NewReader(sortedBookmarkFile))
	if err != nil {
		t.Fatal(err)
	}
	err = bookmarks.Write(&buf)
	if err != nil {
		t.Fatal(err)
	}

	if sortedBookmarkFile != buf.String() {
		t.Fail()
	}
}

func TestReadWriteUnsorted(t *testing.T) {
	var buf bytes.Buffer
	bookmarks := bookmarks.New()

	err := bookmarks.Read(strings.NewReader(unsortedBookmarkFile))
	if err != nil {
		t.Fatal(err)
	}
	err = bookmarks.Write(&buf)
	if err != nil {
		t.Fatal(err)
	}

	if sortedBookmarkFile != buf.String() {
		t.Fail()
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

	for i := 0; i < 10; i++ {
		bookmarks.Add(i)
	}
	for i := 0; i < 10; i++ {
		bookmarks.Remove(i)
	}

	close(ch)
	<-done

	if notifyCount != 20 {
		t.Fail()
	}
}

func TestRemoveObserver(t *testing.T) {
	ch := make(chan string)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ch:
				t.Fail()
			case <-done:
				return
			}
		}
	}()

	bookmarks := bookmarks.New()
	bookmarks.RemoveObserver(bookmarks.AddObserver(ch))

	for i := 0; i < 10; i++ {
		bookmarks.Add(i)
	}

	done <- struct{}{}
}
