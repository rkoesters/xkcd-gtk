package widget

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
)

const (
	// ImageScaleMin is the minimum scale by which the application will zoom
	// the comic image.
	ImageScaleMin = 0.25
	// ImageScaleMax is the maximum scale by which the application will zoom
	// the comic image.
	ImageScaleMax = 5
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

	ImageScale float64

	PropertiesVisible   bool
	PropertiesHeight    int
	PropertiesWidth     int
	PropertiesPositionX int
	PropertiesPositionY int
}

func (ws *WindowState) loadDefaults() {
	newestComic, _ := cache.NewestComicInfoFromCache()
	ws.ComicNumber = newestComic.Num
	ws.Maximized = false
	ws.Height = 500
	ws.Width = 700
	ws.PositionX = 0
	ws.PositionY = 0
	ws.ImageScale = 1
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
	if ws.ImageScale < ImageScaleMin || ws.ImageScale > ImageScaleMax {
		ws.ImageScale = 1
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
	checkForMisplacedWindowState()

	ws.ReadFile(windowStatePath())
}

type StateHaver interface {
	// IsNil returns whether the value contained in the interface is nil,
	// regardless of whether the type is nil. Useful for interfaces because
	// interfaces are a tuple of (type, value), and (type, nil) != nil, only
	// (nil, nil) == nil.
	IsNil() bool

	IsMaximized() bool
	GetSize() (int, int)     // returns width, height
	GetPosition() (int, int) // returns x, y
}

// SaveState writes win.state to disk so it can be loaded next time we open a
// window.
func (ws *WindowState) SaveState(window, dialog StateHaver) {
	ws.Maximized = window.IsMaximized()
	if !ws.Maximized {
		ws.Width, ws.Height = window.GetSize()
		ws.PositionX, ws.PositionY = window.GetPosition()
	}
	ws.PropertiesVisible = !dialog.IsNil()
	if ws.PropertiesVisible {
		ws.PropertiesWidth, ws.PropertiesHeight = dialog.GetSize()
		ws.PropertiesPositionX, ws.PropertiesPositionY = dialog.GetPosition()
	}

	err := ws.WriteFile(windowStatePath())
	if err != nil {
		log.Printf("error saving window state: %v", err)
	}
}

func checkForMisplacedWindowState() {
	misplacedWindowState := filepath.Join(paths.Builder{}.CacheDir(), "state")

	_, err := os.Stat(misplacedWindowState)
	if !os.IsNotExist(err) {
		log.Printf("WARNING: Potentially misplaced window state file '%v'. Should be '%v'.", misplacedWindowState, windowStatePath())
	}
}

func windowStatePath() string {
	return filepath.Join(paths.CacheDir(), "state")
}
