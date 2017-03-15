package main

import (
	"encoding/json"
	"fmt"
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
}

func NewWindowState(w *Window) *WindowState {
	ws := new(WindowState)
	ws.ComicNumber = w.comic.Num
	ws.Width, ws.Height = w.win.GetSize()
	ws.PositionX, ws.PositionY = w.win.GetPosition()
	return ws
}

func (ws *WindowState) String() string {
	return fmt.Sprintf("WindowState{ ComicNumber: %v, Height: %v, Width: %v PositionX: %v, PositionY: %v }", ws.ComicNumber, ws.Height, ws.Width, ws.PositionX, ws.PositionY)
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
	log.Printf("writing state to %v", filename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return ws.Write(f)
}
