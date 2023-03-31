// Package style provides custom application CSS as well as other styling
// utilities.
package style

import (
	_ "embed"
	"sync"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/log"
)

const (
	ClassComicContainer      = "comic-container"
	ClassLinked              = "linked"
	ClassNoMinWidth          = "no-min-width"
	ClassSlimButton          = "slim-button"
	ClassFixHiddenComicTitle = "fix-hidden-comic-title"

	PaddingComicListButton   = 8
	PaddingPopover           = 10
	PaddingPopoverCompact    = 8
	PaddingAuxiliaryWindow   = 12
	PaddingUnlinkedButtonBox = 4
)

//go:embed light.css
var lightCSS string

//go:embed dark.css
var darkCSS string

var (
	cssDataMutex      sync.RWMutex
	cssProvider       *gtk.CssProvider // Protected by cssDataMutex
	loadedCSSDarkMode bool             // Protected by cssDataMutex
)

// InitCSS initializes the application's custom CSS.
func InitCSS(darkMode bool) error {
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

	return loadCSS(cssProvider, darkMode)
}

// UpdateCSS reloads the application CSS if darkMode does not match the
// currently loaded CSS.
func UpdateCSS(darkMode bool) error {
	log.Debugf("UpdateCSS(darkMode=%v)", darkMode)
	cssDataMutex.RLock()
	if darkMode == loadedCSSDarkMode {
		cssDataMutex.RUnlock()
		return nil
	}
	cssDataMutex.RUnlock()

	cssDataMutex.Lock()
	defer cssDataMutex.Unlock()

	return loadCSS(cssProvider, darkMode)
}

func loadCSS(p *gtk.CssProvider, darkMode bool) error {
	loadedCSSDarkMode = darkMode
	if darkMode {
		log.Debug("loading dark.css")
		return p.LoadFromData(darkCSS)
	} else {
		log.Debug("loading light.css")
		return p.LoadFromData(lightCSS)
	}
}
