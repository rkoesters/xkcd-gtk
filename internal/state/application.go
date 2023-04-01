// Package state provides data structures for storing user state.
package state

import (
	"encoding/json"
	"io"
	"os"
)

// Application is a struct that holds our application's settings.
type Application struct {
	DarkMode bool
}

func (a *Application) loadDefaults() {
	a.DarkMode = false
}

// Read takes the given io.Reader and tries to parse json encoded state from it.
func (a *Application) Read(r io.Reader) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(a)
	if err != nil {
		a.loadDefaults()
		return err
	}
	return nil
}

// ReadFile opens the given file and calls Read on the contents.
func (a *Application) ReadFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		a.loadDefaults()
		return err
	}
	defer f.Close()
	return a.Read(f)
}

// Write takes the given io.Writer and writes the Settings struct to it in json.
func (a *Application) WriteTo(w io.Writer) (int64, error) {
	bc := &byteCounter{Writer: w}
	enc := json.NewEncoder(bc)
	err := enc.Encode(a)
	return bc.count, err
}

// WriteFile creates or truncates the given file and calls Write on it.
func (a *Application) WriteFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = a.WriteTo(f)
	return err
}
