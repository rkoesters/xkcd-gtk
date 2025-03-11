package state_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/rkoesters/xkcd-gtk/internal/state"
)

const (
	state1 = `{"DarkMode":true}
`
	state2 = `{"DarkMode":false}
`
)

func TestReadFrom(t *testing.T) {
	var as state.Application

	r := strings.NewReader(state1)
	as.ReadFrom(r)
	if !as.DarkMode {
		t.Fatal("dark mode is disabled, config: ", state1)
	}

	r = strings.NewReader(state2)
	as.ReadFrom(r)
	if as.DarkMode {
		t.Fatal("dark mode is enabled, config: ", state1)
	}
}

func TestWrite(t *testing.T) {
	var as state.Application
	var b strings.Builder

	as.DarkMode = true
	as.WriteTo(&b)
	if b.String() != state1 {
		t.Fatalf("b.String()=%q state1=%q", b.String(), state1)
	}

	b.Reset()

	as.DarkMode = false
	as.WriteTo(&b)
	if b.String() != state2 {
		t.Fatalf("b.String()=%q state2=%q", b.String(), state2)
	}
}

func TestBadReadFrom(t *testing.T) {
	var as state.Application

	r := strings.NewReader("bad format")
	as.ReadFrom(r)
	if as.DarkMode {
		t.Fatal("dark mode enabled after bad read")
	}
}

func TestApplicationState(t *testing.T) {
	tests := []struct {
		name string
		app  state.Application
	}{{
		name: "dark on",
		app:  state.Application{DarkMode: true},
	}, {
		name: "dark off",
		app:  state.Application{DarkMode: false},
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(),
				strings.ReplaceAll(test.name, " ", "_"))

			err := test.app.WriteFile(path)
			if err != nil {
				t.Fatalf("error writing %q: %v", path, err)
			}

			var app state.Application
			err = app.ReadFile(path)
			if err != nil {
				t.Fatalf("error reading %q: %v", path, err)
			}

			if test.app.DarkMode != app.DarkMode {
				t.Error("mismatch between WriteFile and ReadFile")
			}
		})
	}
}
