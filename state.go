package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// WindowState is a struct that holds the information about the state of
// a Window. This struct is meant to be stored so we can restore the
// state of a Window.
type WindowState struct {
	ComicNumber int
	Height      int
	Width       int
	PositionX   int
	PositionY   int
	Maximized   bool

	PropertiesVisible   bool
	PropertiesHeight    int
	PropertiesWidth     int
	PropertiesPositionX int
	PropertiesPositionY int
}

// NewWindowState creates a WindowState struct by looking at the state
// of the given Window.
func NewWindowState(w *Window) *WindowState {
	ws := new(WindowState)
	ws.ComicNumber = w.comic.Num
	ws.Width, ws.Height = w.win.GetSize()
	ws.PositionX, ws.PositionY = w.win.GetPosition()
	ws.Maximized = w.win.IsMaximized()
	if w.properties != nil {
		ws.PropertiesVisible = true
		ws.PropertiesWidth, ws.PropertiesHeight = w.properties.dialog.GetSize()
		ws.PropertiesPositionX, ws.PropertiesPositionY = w.properties.dialog.GetPosition()
	}
	return ws
}

// Read takes the given io.Reader and tries to parse json encoded state
// from it.
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

// ReadFile opens the given file and calls Read on the contents.
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

// Write takes the given io.Writer and writes the WindowState struct to
// it in json.
func (ws *WindowState) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(ws)
}

// WriteFile creates or truncates the given file and calls Write on it.
func (ws *WindowState) WriteFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("writing state to %v: %v", filename, err)
		return err
	}
	defer f.Close()
	return ws.Write(f)
}
