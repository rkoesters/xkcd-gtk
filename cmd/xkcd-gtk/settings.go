package main

import (
	"encoding/json"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Settings is a struct that holds our application's settings.
type Settings struct {
	DarkMode bool
}

func (settings *Settings) loadDefaults() {
	settings.DarkMode = false
}

// Read takes the given io.Reader and tries to parse json encoded state from it.
func (settings *Settings) Read(r io.Reader) {
	dec := json.NewDecoder(r)
	err := dec.Decode(settings)
	if err != nil {
		settings.loadDefaults()
	}
}

// ReadFile opens the given file and calls Read on the contents.
func (settings *Settings) ReadFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		settings.loadDefaults()
		return
	}
	defer f.Close()
	settings.Read(f)
}

// Write takes the given io.Writer and writes the Settings struct to it in json.
func (settings *Settings) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(settings)
}

// WriteFile creates or truncates the given file and calls Write on it.
func (settings *Settings) WriteFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return settings.Write(f)
}

// LoadSettings tries to load our settings from disk.
func (app *Application) LoadSettings() {
	var err error

	// Read settings from disk.
	app.settings.ReadFile(settingsPath())

	// Get reference to Gtk's settings.
	app.gtkSettings, err = gtk.SettingsGetDefault()
	if err == nil {
		// Apply Dark Mode setting.
		err = app.gtkSettings.SetProperty("gtk-application-prefer-dark-theme", app.settings.DarkMode)
		if err != nil {
			log.Print("error setting dark mode state: ", err)
		}
	} else {
		log.Print("error querying gtk settings: ", err)
	}
}

// SaveSettings tries to save our settings to disk.
func (app *Application) SaveSettings() {
	err := os.MkdirAll(paths.ConfigDir(), 0755)
	if err != nil {
		log.Printf("error saving settings: %v", err)
	}

	err = app.settings.WriteFile(settingsPath())
	if err != nil {
		log.Printf("error saving settings: %v", err)
	}
}

func settingsPath() string {
	return filepath.Join(paths.ConfigDir(), "settings")
}
