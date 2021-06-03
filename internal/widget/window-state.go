package widget

import (
	"encoding/json"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"io"
	"log"
	"os"
	"path/filepath"
)

// WindowState is a struct that holds the information about the state of a
// Window. This struct is meant to be stored so we can restore the state of a
// Window.
type WindowState struct {
	ComicNumber int

	Maximized bool
	Height    int
	Width     int
	PositionX int
	PositionY int

	PropertiesVisible   bool
	PropertiesHeight    int
	PropertiesWidth     int
	PropertiesPositionX int
	PropertiesPositionY int
}

func (ws *WindowState) loadDefaults() {
	newestComic, _ := cache.NewestComicInfo()
	ws.ComicNumber = newestComic.Num
	ws.Maximized = false
	ws.Height = 500
	ws.Width = 700
	ws.PositionX = 0
	ws.PositionY = 0
	ws.PropertiesVisible = false
	ws.PropertiesHeight = 350
	ws.PropertiesWidth = 300
	ws.PropertiesPositionX = 0
	ws.PropertiesPositionY = 0
}

// Read takes the given io.Reader and tries to parse json encoded state from it.
func (ws *WindowState) Read(r io.Reader) {
	dec := json.NewDecoder(r)
	err := dec.Decode(ws)
	if err != nil {
		ws.loadDefaults()
	}
}

// ReadFile opens the given file and calls Read on the contents.
func (ws *WindowState) ReadFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		ws.loadDefaults()
		return
	}
	defer f.Close()
	ws.Read(f)
}

// Write takes the given io.Writer and writes the WindowState struct to it in
// json.
func (ws *WindowState) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(ws)
}

// WriteFile creates or truncates the given file and calls Write on it.
func (ws *WindowState) WriteFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return ws.Write(f)
}

func (ws *WindowState) LoadState() {
	ws.ReadFile(windowStatePath())
}

// SaveState writes win.state to disk so it can be loaded next time we open a
// window.
func (ws *WindowState) SaveState(window *gtk.ApplicationWindow, propertiesDialog *gtk.Dialog) {
	ws.Maximized = window.IsMaximized()
	if !ws.Maximized {
		ws.Width, ws.Height = window.GetSize()
		ws.PositionX, ws.PositionY = window.GetPosition()
	}
	if propertiesDialog == nil {
		ws.PropertiesVisible = false
	} else {
		ws.PropertiesVisible = true
		ws.PropertiesWidth, ws.PropertiesHeight = propertiesDialog.GetSize()
		ws.PropertiesPositionX, ws.PropertiesPositionY = propertiesDialog.GetPosition()
	}

	err := ws.WriteFile(windowStatePath())
	if err != nil {
		log.Printf("error saving window state: %v", err)
	}
}

func windowStatePath() string {
	return filepath.Join(paths.CacheDir(), "state")
}
