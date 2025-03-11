package state

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
	// ImageScaleMin is the minimum scale by which the application will zoom the
	// comic image.
	ImageScaleMin = 0.25
	// ImageScaleMax is the maximum scale by which the application will zoom the
	// comic image.
	ImageScaleMax = 5
)

// Window is a struct that holds the information about the state of a Window.
// This struct is meant to be stored so we can restore the state of a Window.
type Window struct {
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

var (
	_ io.WriterTo   = &Window{}
	_ io.ReaderFrom = &Window{}
)

func (w *Window) loadDefaults() {
	newestComic, _ := cache.NewestComicInfoFromCache()
	w.ComicNumber = newestComic.Num
	w.Maximized = false
	w.Height = 500
	w.Width = 700
	w.PositionX = 0
	w.PositionY = 0
	w.ImageScale = 1
	w.PropertiesVisible = false
	w.PropertiesHeight = 350
	w.PropertiesWidth = 300
	w.PropertiesPositionX = 0
	w.PropertiesPositionY = 0
}

func (w *Window) HasPosition() bool {
	return w.PositionX != 0 && w.PositionY != 0
}

func (w *Window) HasPropertiesPosition() bool {
	return w.PropertiesPositionX != 0 && w.PropertiesPositionY != 0
}

// ReadFrom takes the given io.Reader and tries to parse json encoded state from
// it.
func (w *Window) ReadFrom(r io.Reader) (int64, error) {
	bc := &byteCounter{Reader: r}
	dec := json.NewDecoder(bc)
	err := dec.Decode(w)
	if err != nil {
		w.loadDefaults()
	}
	if w.ImageScale < ImageScaleMin || w.ImageScale > ImageScaleMax {
		w.ImageScale = 1
	}
	return bc.count, err
}

// ReadFile opens the given file and calls Read on the contents.
func (w *Window) ReadFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		w.loadDefaults()
		return
	}
	defer f.Close()
	w.ReadFrom(f)
}

// WriteTo takes the given io.Writer and writes the Window struct to it in json.
func (w *Window) WriteTo(o io.Writer) (int64, error) {
	bc := &byteCounter{Writer: o}
	enc := json.NewEncoder(bc)
	err := enc.Encode(w)
	return bc.count, err
}

// WriteFile creates or truncates the given file and calls Write on it.
func (w *Window) WriteFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = w.WriteTo(f)
	return err
}

func (w *Window) LoadState() {
	checkForMisplacedWindowState()

	w.ReadFile(windowStatePath())
}

type StateHaver interface {
	IsVisible() bool
	IsMaximized() bool
	GetSize() (int, int)     // returns width, height
	GetPosition() (int, int) // returns x, y
}

// SaveState writes win.state to disk so it can be loaded next time we open a
// window.
func (w *Window) SaveState(window, dialog StateHaver) {
	w.Maximized = window.IsMaximized()
	w.Width, w.Height = window.GetSize()
	w.PositionX, w.PositionY = window.GetPosition()

	w.PropertiesVisible = dialog.IsVisible()
	if w.PropertiesVisible {
		w.PropertiesWidth, w.PropertiesHeight = dialog.GetSize()
		w.PropertiesPositionX, w.PropertiesPositionY = dialog.GetPosition()
	}

	err := w.WriteFile(windowStatePath())
	if err != nil {
		log.Printf("error saving window state: %v", err)
	}
}

func checkForMisplacedWindowState() {
	misplacedWindowState := filepath.Join(paths.Builder{}.CacheDir(), "state")

	_, err := os.Stat(misplacedWindowState)
	if !os.IsNotExist(err) {
		log.Printf("WARNING: Potentially misplaced window state file %q. Should be %q.", misplacedWindowState, windowStatePath())
	}
}

func windowStatePath() string {
	return filepath.Join(paths.CacheDir(), "state")
}
