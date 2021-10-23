package bookmarks

import (
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"log"
	"os"
	"path/filepath"
)

// Path returns the path to the user's bookmarks file.
func Path() string {
	return filepath.Join(paths.DataDir(), "bookmarks")
}

// CheckForMisplacedFiles prints a warning message to standard error if there
// are any stray bookmark files that may have been caused by a bug that commit
// d13e4dc0ff81e9d12df29e7f9be4e82e7f70cc01 fixed.
func CheckForMisplacedFiles() {
	misplacedBookmarksList := []string{
		filepath.Join(paths.Builder{}.ConfigDir(), "bookmarks"),
		filepath.Join(paths.Builder{}.DataDir(), "bookmarks"),
		filepath.Join(paths.ConfigDir(), "bookmarks"),
	}

	for _, p := range misplacedBookmarksList {
		_, err := os.Stat(p)
		if !os.IsNotExist(err) {
			log.Printf("WARNING: Potentially misplaced bookmarks file '%v'. Should be '%v'.", p, Path())
		}
	}
}
