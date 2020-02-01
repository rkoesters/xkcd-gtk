package paths

import (
	"github.com/rkoesters/xdg/basedir"
	"os"
	"path/filepath"
)

type builder struct {
	appID string
}

// CacheDir returns the path to our app's cache directory.
func (b builder) CacheDir() string {
	return filepath.Join(basedir.CacheHome, b.appID)
}

// ConfigDir returns the path to our app's user configuration directory.
func (b builder) ConfigDir() string {
	return filepath.Join(basedir.ConfigHome, b.appID)
}

// DataDir returns the path to our app's user data directory.
func (b builder) DataDir() string {
	return filepath.Join(basedir.DataHome, b.appID)
}

// LocaleDir returns the path to the system locale directory.
func (b builder) LocaleDir() string {
	for _, dir := range basedir.DataDirs {
		path := filepath.Join(dir, "locale")
		_, err := os.Stat(path)
		if err == nil {
			return path
		}
	}
	return "."
}
