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
}

func NewWindowState(w *Window) *WindowState {
	ws := new(WindowState)
	ws.ComicNumber = w.comic.Num
	ws.Height = w.win.GetAllocatedHeight()
	ws.Width = w.win.GetAllocatedWidth()
	log.Printf("Allocation: %v", ws)
	ws.Width, ws.Height = w.win.GetSize()
	log.Printf("GetSize: %v", ws)
	return ws
}

func (ws *WindowState) Read(r io.Reader) {
	log.Print("reading state")
	dec := json.NewDecoder(r)
	err := dec.Decode(ws)
	if err != nil {
		// Something is wrong, lets load defaults.
		newestComic, _ := GetNewestComicInfo()
		ws.ComicNumber = newestComic.Num
		ws.Height = 800
		ws.Width = 1000
	}
}

func (ws *WindowState) ReadFile(filename string) {
	log.Printf("reading state from %v", filename)
	f, err := os.Open(filename)
	if err != nil {
		// Can't read file, lets load defaults.
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
	log.Print("writing state")
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
