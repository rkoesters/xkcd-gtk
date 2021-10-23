package settings

import (
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"log"
	"os"
	"path/filepath"
)

// Path returns the path to the user's settings file.
func Path() string {
	return filepath.Join(paths.ConfigDir(), "settings")
}

// CheckForMisplacedFiles prints a warning message to standard error if there
// are any stray configuration files that may have been caused by a bug that
// commit d13e4dc0ff81e9d12df29e7f9be4e82e7f70cc01 fixed.
func CheckForMisplacedFiles() {
	misplacedSettings := filepath.Join(paths.Builder{}.ConfigDir(), "settings")

	_, err := os.Stat(misplacedSettings)
	if !os.IsNotExist(err) {
		log.Printf("WARNING: Potentially misplaced settings file '%v'. Should be '%v'.", misplacedSettings, Path())
	}
}
