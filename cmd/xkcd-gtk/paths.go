package main

import (
	"github.com/rkoesters/xdg/basedir"
	"os"
	"path/filepath"
)

// CacheDir returns the path to our app's cache directory.
func CacheDir() string {
	return filepath.Join(basedir.CacheHome, appID)
}

// ConfigDir returns the path to our app's user configuration directory.
func ConfigDir() string {
	return filepath.Join(basedir.ConfigHome, appID)
}

// DataDir returns the path to our app's user data directory.
func DataDir() string {
	return filepath.Join(basedir.DataHome, appID)
}

// LocaleDir returns the path to the system locale directory.
func LocaleDir() string {
	for _, dir := range basedir.DataDirs {
		path := filepath.Join(dir, "locale")
		_, err := os.Stat(path)
		if err == nil {
			return path
		}
	}
	return "."
}
