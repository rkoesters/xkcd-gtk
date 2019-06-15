package main

import (
	"github.com/rkoesters/xdg"
	"log"
)

func openURL(url string) {
	err := xdg.Open(url)
	if err != nil {
		log.Print("error opening ", url, " in web browser: ", err)
	}
}
