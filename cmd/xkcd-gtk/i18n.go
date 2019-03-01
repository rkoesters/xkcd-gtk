package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/rkoesters/xdg/basedir"
	"os"
	"path/filepath"
)

var l = glib.Local

func init() {
	glib.InitI18n(appID, localeDir())
}

func localeDir() string {
	for _, dir := range basedir.DataDirs {
		path := filepath.Join(dir, "locale")
		_, err := os.Stat(path)
		if err == nil {
			return path
		}
	}
	return "."
}
