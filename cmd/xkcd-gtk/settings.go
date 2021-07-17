package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"log"
	"os"
	"path/filepath"
)

// LoadSettings tries to load our settings from disk.
func (app *Application) LoadSettings() {
	var err error

	checkForMisplacedSettings()

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
	err := paths.EnsureConfigDir()
	if err != nil {
		log.Printf("error saving settings: %v", err)
	}

	err = app.settings.WriteFile(settingsPath())
	if err != nil {
		log.Printf("error saving settings: %v", err)
	}
}

func checkForMisplacedSettings() {
	misplacedSettings := filepath.Join(paths.Builder{}.ConfigDir(), "settings")

	_, err := os.Stat(misplacedSettings)
	if !os.IsNotExist(err) {
		log.Printf("WARNING: Potentially misplaced settings file '%v'. Should be '%v'.", misplacedSettings, settingsPath())
	}
}

func settingsPath() string {
	return filepath.Join(paths.ConfigDir(), "settings")
}
