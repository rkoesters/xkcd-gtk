package paths

import (
	"os"
	"path/filepath"

	"github.com/rkoesters/xkcd-gtk/internal/log"
)

// SearchIndex returns the path to the application's search index.
func (b Builder) SearchIndex() string {
	return filepath.Join(b.CacheDir(), "search")
}

// SearchIndex returns the path to the application's search index.
func SearchIndex() string {
	return b.SearchIndex()
}

// CheckForMisplacedSearchIndex prints a warning message to standard error if
// there are any stray bookmark files that may have been caused by a bug that
// commit d13e4dc0ff81e9d12df29e7f9be4e82e7f70cc01 fixed.
func CheckForMisplacedSearchIndex() {
	misplacedSearchIndex := Builder{}.SearchIndex()

	_, err := os.Stat(misplacedSearchIndex)
	if !os.IsNotExist(err) {
		log.Printf("WARNING: Potentially misplaced search index %q. Should be %q.", misplacedSearchIndex, SearchIndex())
	}
}
