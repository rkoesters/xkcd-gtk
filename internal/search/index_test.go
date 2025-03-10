package search_test

import (
	"path/filepath"
	"strconv"
	"testing"

	"github.com/rkoesters/xkcd"
	"github.com/rkoesters/xkcd-gtk/internal/search"
)

const (
	testComicNumber = 404
	testComicTitle  = "test comic"
)

func TestSearchIndex(t *testing.T) {
	path := filepath.Join(t.TempDir(), "search")

	si, err := search.New(path)
	if err != nil {
		t.Fatalf("error creating test search index %q: %v", path, err)
	}

	comic := &xkcd.Comic{
		Num:   testComicNumber,
		Title: testComicTitle,
	}
	err = si.Index(comic)
	if err != nil {
		t.Errorf("error indexing comic %q: %v", comic, err)
	}

	results, err := si.Search(testComicTitle)
	if err != nil {
		si.Close()
		t.Fatalf("error searching index with query %q: %v", testComicTitle, err)
	}

	if results.Total != 1 {
		t.Errorf("expected 1 result, got %q", results.Total)
	}

	for _, result := range results.Hits {
		n, err := strconv.Atoi(result.ID)
		if err != nil {
			t.Error("error converting search result key into integer: ", err)
			continue
		}
		if n != testComicNumber {
			t.Error("unexpected result: ", n)
		}
	}

	err = si.Close()
	if err != nil {
		t.Error("error closing search index: ", err)
	}
}
