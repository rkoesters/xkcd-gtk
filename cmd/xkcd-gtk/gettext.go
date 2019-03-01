package main

import (
	"github.com/leonelquinteros/gotext"
	"github.com/rkoesters/xdg/basedir"
	"github.com/rkoesters/xdg/keyfile"
	"os"
	"path/filepath"
)

var (
	gt  = gotext.Get
	gtn = gotext.GetN
)

func init() {
	gotext.Configure(localeDir(), defaultLocale(), appID)
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

func defaultLocale() string {
	return keyfile.DefaultLocale().String()
}
