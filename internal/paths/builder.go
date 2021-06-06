package paths

import (
	"github.com/rkoesters/xdg/basedir"
	"os"
	"path/filepath"
)

const defaultDirMode = 0755

type builder struct {
	appID string
}

// CacheDir returns the path to our app's cache directory.
func (b builder) CacheDir() string {
	return filepath.Join(basedir.CacheHome, b.appID)
}

// EnsureCacheDir creates the CacheDir() directory, if it doesn't exist. It does
// not return an error if it already exists.
func (b builder) EnsureCacheDir() error {
	return os.MkdirAll(b.CacheDir(), defaultDirMode)
}

// ConfigDir returns the path to our app's user configuration directory.
func (b builder) ConfigDir() string {
	return filepath.Join(basedir.ConfigHome, b.appID)
}

// EnsureConfigDir creates the ConfigDir() directory, if it doesn't exist. It
// does not return an error if it already exists.
func (b builder) EnsureConfigDir() error {
	return os.MkdirAll(b.ConfigDir(), defaultDirMode)
}

// DataDir returns the path to our app's user data directory.
func (b builder) DataDir() string {
	return filepath.Join(basedir.DataHome, b.appID)
}

// EnsureDataDir creates the DataDir() directory, if it doesn't exist. It does
// not return an error if it already exists.
func (b builder) EnsureDataDir() error {
	return os.MkdirAll(b.DataDir(), defaultDirMode)
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
