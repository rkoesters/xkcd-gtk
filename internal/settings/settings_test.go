package settings_test

import (
	"strings"
	"testing"

	"github.com/rkoesters/xkcd-gtk/internal/settings"
)

const (
	settings1 = `{"DarkMode":true}
`
	settings2 = `{"DarkMode":false}
`
)

func TestRead(t *testing.T) {
	var settings settings.Settings

	r := strings.NewReader(settings1)
	settings.Read(r)
	if !settings.DarkMode {
		t.Fatal("dark mode is disabled, config: ", settings1)
	}

	r = strings.NewReader(settings2)
	settings.Read(r)
	if settings.DarkMode {
		t.Fatal("dark mode is enabled, config: ", settings1)
	}
}

func TestWrite(t *testing.T) {
	var settings settings.Settings
	var b strings.Builder

	settings.DarkMode = true
	settings.Write(&b)
	if b.String() != settings1 {
		t.Fatalf("b.String()='%v' settings1='%v'", b.String(), settings1)
	}

	b.Reset()

	settings.DarkMode = false
	settings.Write(&b)
	if b.String() != settings2 {
		t.Fatalf("b.String()='%v' settings2='%v'", b.String(), settings2)
	}
}

func TestBadRead(t *testing.T) {
	var settings settings.Settings

	r := strings.NewReader("bad format")
	settings.Read(r)
	if settings.DarkMode {
		t.Fatal("dark mode enabled after bad read")
	}
}
