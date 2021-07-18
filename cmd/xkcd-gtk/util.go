package main

import (
	"github.com/rkoesters/xdg"
	"github.com/rkoesters/xkcd-gtk/internal/log"
)

func openURL(url string) {
	err := xdg.Open(url)
	if err != nil {
		log.Print("error opening ", url, " in web browser: ", err)
	}
}
