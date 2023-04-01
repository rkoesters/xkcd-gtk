package state_test

import (
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
		t.Fatalf("b.String()='%v' state1='%v'", b.String(), state1)
	}

	b.Reset()

	as.DarkMode = false
	as.WriteTo(&b)
	if b.String() != state2 {
		t.Fatalf("b.String()='%v' state2='%v'", b.String(), state2)
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
