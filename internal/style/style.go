// Package style provides custom application CSS as well as other styling
// utilities.
package style

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"regexp"
	"strings"
)

const (
	ClassComicContainer = "comic-container"
	ClassDark           = "dark"
	ClassLinked         = "linked"
)

// LoadCSS provides the application's custom CSS to GTK.
func LoadCSS() error {
	provider, err := gtk.CssProviderNew()
	if err != nil {
		return err
	}
	provider.LoadFromData(styleCSS)

	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		return err
	}

	gtk.AddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return nil
}

var (
	// largeToolbarThemes is the list of gtk themes for which we should use
	// large toolbar buttons.
	largeToolbarThemesRegexp = regexp.MustCompile(strings.Join([]string{
		"elementary(-x)?",
		"io\\.elementary\\.stylesheet.*",
		"win32",
	}, "|"))

	// nonSymbolicIconThemes is the list of gtk themes for which we should
	// use non-symbolic icons.
	nonSymbolicIconThemesRegexp = regexp.MustCompile(strings.Join([]string{
		"elementary(-x)?",
		"io\\.elementary\\.stylesheet.*",
	}, "|"))
)

func IsLargeToolbarTheme(theme string) bool {
	return largeToolbarThemesRegexp.MatchString(theme)
}

func IsSymbolicIconTheme(theme string) bool {
	return !nonSymbolicIconThemesRegexp.MatchString(theme)
}
