package paths

import (
	"os"
	"path/filepath"

	"github.com/rkoesters/xkcd-gtk/internal/log"
)

// Settings returns the path to the user's settings file.
func (b Builder) Settings() string {
	return filepath.Join(b.ConfigDir(), "settings")
}

// Settings returns the path to the user's settings file.
func Settings() string {
	return b.Settings()
}

// CheckForMisplacedSettings prints a warning message to standard error if there
// are any stray configuration files that may have been caused by a bug that
// commit d13e4dc0ff81e9d12df29e7f9be4e82e7f70cc01 fixed.
func CheckForMisplacedSettings() {
	misplacedSettings := Builder{}.Settings()

	_, err := os.Stat(misplacedSettings)
	if !os.IsNotExist(err) {
		log.Printf("WARNING: Potentially misplaced settings file %q. Should be %q.", misplacedSettings, Settings())
	}
}
