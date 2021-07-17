// Package settings provides data structures for storing user settings.
package settings

import (
	"encoding/json"
	"io"
	"os"
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
