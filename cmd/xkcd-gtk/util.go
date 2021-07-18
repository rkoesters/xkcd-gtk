package main

import (
	"github.com/rkoesters/xdg"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"os"
	"path/filepath"
)

func openURL(url string) {
	err := xdg.Open(url)
	if err != nil {
		log.Print("error opening ", url, " in web browser: ", err)
	}
}

func checkForMisplacedBookmarks() {
	misplacedBookmarksList := []string{
		filepath.Join(paths.Builder{}.ConfigDir(), "bookmarks"),
		filepath.Join(paths.Builder{}.DataDir(), "bookmarks"),
		filepath.Join(paths.ConfigDir(), "bookmarks"),
	}

	for _, p := range misplacedBookmarksList {
		_, err := os.Stat(p)
		if !os.IsNotExist(err) {
			log.Printf("WARNING: Potentially misplaced bookmarks file '%v'. Should be '%v'.", p, bookmarksPath())
		}
	}
}

func checkForMisplacedSettings() {
	misplacedSettings := filepath.Join(paths.Builder{}.ConfigDir(), "settings")

	_, err := os.Stat(misplacedSettings)
	if !os.IsNotExist(err) {
		log.Printf("WARNING: Potentially misplaced settings file '%v'. Should be '%v'.", misplacedSettings, settingsPath())
	}
}

func bookmarksPath() string {
	return filepath.Join(paths.DataDir(), "bookmarks")
}

func settingsPath() string {
	return filepath.Join(paths.ConfigDir(), "settings")
}
