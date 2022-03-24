package paths

import (
	"log"
	"os"
	"path/filepath"
)

// Bookmarks returns the path to the user's bookmarks file.
func (b Builder) Bookmarks() string {
	return filepath.Join(b.DataDir(), "bookmarks")
}

// Bookmarks returns the path to the user's bookmarks file.
func Bookmarks() string {
	return b.Bookmarks()
}

// CheckForMisplacedBookmarks prints a warning message to standard error if
// there are any stray bookmark files that may have been caused by a bug that
// commit d13e4dc0ff81e9d12df29e7f9be4e82e7f70cc01 fixed.
func CheckForMisplacedBookmarks() {
	misplacedBookmarksList := []string{
		filepath.Join(Builder{}.ConfigDir(), "bookmarks"),
		filepath.Join(Builder{}.Bookmarks()),
		filepath.Join(ConfigDir(), "bookmarks"),
	}

	for _, p := range misplacedBookmarksList {
		_, err := os.Stat(p)
		if !os.IsNotExist(err) {
			log.Printf("WARNING: Potentially misplaced bookmarks file '%v'. Should be '%v'.", p, Bookmarks())
		}
	}
}