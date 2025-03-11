package paths

import (
	"os"
	"path/filepath"

	"github.com/rkoesters/xdg/basedir"
	"github.com/rkoesters/xkcd-gtk/internal/log"
)

const defaultDirMode = 0755

// Builder provides methods to find the paths where the app should store the
// files it creates.
type Builder struct {
	appID string
}

// CacheDir returns the path to our app's cache directory.
func (b Builder) CacheDir() string {
	return filepath.Join(basedir.CacheHome, b.appID)
}

// EnsureCacheDir creates the CacheDir() directory, if it doesn't exist. It does
// not return an error if it already exists.
func (b Builder) EnsureCacheDir() error {
	p := b.CacheDir()
	log.Debugf("Ensuring cache directory %q exists", p)
	return os.MkdirAll(p, defaultDirMode)
}

// ConfigDir returns the path to our app's user configuration directory.
func (b Builder) ConfigDir() string {
	return filepath.Join(basedir.ConfigHome, b.appID)
}

// EnsureConfigDir creates the ConfigDir() directory, if it doesn't exist. It
// does not return an error if it already exists.
func (b Builder) EnsureConfigDir() error {
	p := b.ConfigDir()
	log.Debugf("Ensuring configuration directory %q exists", p)
	return os.MkdirAll(p, defaultDirMode)
}

// DataDir returns the path to our app's user data directory.
func (b Builder) DataDir() string {
	return filepath.Join(basedir.DataHome, b.appID)
}

// EnsureDataDir creates the DataDir() directory, if it doesn't exist. It does
// not return an error if it already exists.
func (b Builder) EnsureDataDir() error {
	p := b.DataDir()
	log.Debugf("Ensuring data directory %q exists", p)
	return os.MkdirAll(p, defaultDirMode)
}

// LocaleDir returns the path to the system locale directory.
func (b Builder) LocaleDir() string {
	for _, dir := range basedir.DataDirs {
		path := filepath.Join(dir, "locale")
		_, err := os.Stat(path)
		if err == nil {
			return path
		}
	}
	return "."
}
