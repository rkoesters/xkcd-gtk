// Package style provides custom application CSS as well as other styling
// utilities.
package style

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"regexp"
	"strings"
	"sync"
)

const (
	ClassComicContainer = "comic-container"
)

var (
	cssDataMutex      sync.RWMutex
	cssProvider       *gtk.CssProvider // Protected by cssDataMutex
	loadedCSSDarkMode bool             // Protected by cssDataMutex
)

// InitCSS initializes the application's custom CSS.
func InitCSS() error {
	var err error

	cssDataMutex.Lock()
	defer cssDataMutex.Unlock()

	cssProvider, err = gtk.CssProviderNew()
	if err != nil {
		return err
	}

	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		return err
	}

	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	// loadedCSSDarkMode defaults to false, and InitCSS is usually called
	// before it is set to anything else.
	return loadCSS(cssProvider, loadedCSSDarkMode)
}

// UpdateCSS reloads the application CSS if darkMode does not match the
// currently loaded CSS.
func UpdateCSS(darkMode bool) error {
	cssDataMutex.RLock()
	if darkMode == loadedCSSDarkMode {
		cssDataMutex.RUnlock()
		return nil
	}
	cssDataMutex.RUnlock()

	cssDataMutex.Lock()
	defer cssDataMutex.Unlock()

	err := loadCSS(cssProvider, darkMode)
	if err != nil {
		return err
	}

	loadedCSSDarkMode = darkMode

	return nil
}

func loadCSS(p *gtk.CssProvider, darkMode bool) error {
	if darkMode {
		log.Debug("loading style-dark.css")
		return p.LoadFromData(styleDarkCSS)
	} else {
		log.Debug("loading style.css")
		return p.LoadFromData(styleCSS)
	}
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

func IsSymbolicIconTheme(theme string, darkMode bool) bool {
	return darkMode || !nonSymbolicIconThemesRegexp.MatchString(theme)
}
