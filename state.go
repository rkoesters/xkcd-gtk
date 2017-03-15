package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type WindowState struct {
	ComicNumber int
	Height      int
	Width       int
	PositionX   int
	PositionY   int

	PropertiesVisible   bool
	PropertiesHeight    int
	PropertiesWidth     int
	PropertiesPositionX int
	PropertiesPositionY int
}

func NewWindowState(w *Window) *WindowState {
	ws := new(WindowState)
	ws.ComicNumber = w.comic.Num
	ws.Width, ws.Height = w.win.GetSize()
	ws.PositionX, ws.PositionY = w.win.GetPosition()
	if w.properties != nil {
		ws.PropertiesVisible = true
		ws.PropertiesWidth, ws.PropertiesHeight = w.properties.dialog.GetSize()
		ws.PropertiesPositionX, ws.PropertiesPositionY = w.properties.dialog.GetPosition()
	}
	return ws
}

func (ws *WindowState) Read(r io.Reader) {
	dec := json.NewDecoder(r)
	err := dec.Decode(ws)
	if err != nil {
		// Something is wrong, lets load defaults.
		log.Printf("reading state: %v", err)
		newestComic, _ := GetNewestComicInfo()
		ws.ComicNumber = newestComic.Num
		ws.Height = 800
		ws.Width = 1000
	}
}

func (ws *WindowState) ReadFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		// Can't read file, lets load defaults.
		log.Printf("reading state from %v: %v", filename, err)
		newestComic, _ := GetNewestComicInfo()
		ws.ComicNumber = newestComic.Num
		ws.Height = 800
		ws.Width = 1000
		return
	}
	defer f.Close()
	ws.Read(f)
}

func (ws *WindowState) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(ws)
}

func (ws *WindowState) WriteFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("writing state to %v: %v", filename, err)
		return err
	}
	defer f.Close()
	return ws.Write(f)
}
