package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	styleClassComicContainer = "comic-container"
	styleClassDark           = "dark"
	styleClassLinked         = "linked"
)

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

// LoadCSS provides the application's custom CSS to GTK.
func (app *Application) LoadCSS() {
	provider, err := gtk.CssProviderNew()
	if err != nil {
		log.Print(err)
		return
	}
	provider.LoadFromData(styleCSS)

	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		log.Print(err)
		return
	}

	gtk.AddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

// StyleUpdated is called when the style of our gtk window is updated.
func (win *Window) StyleUpdated() {
	// First, lets find out what GTK theme we are using.
	themeName := os.Getenv("GTK_THEME")
	if themeName == "" {
		// The theme is not being set by the environment, so lets ask
		// GTK what theme it is going to use.
		themeNameIface, err := win.app.gtkSettings.GetProperty("gtk-theme-name")
		if err != nil {
			log.Print(err)
		} else {
			themeNameStr, ok := themeNameIface.(string)
			if ok {
				themeName = themeNameStr
			}
		}
	}

	// The default size for our headerbar buttons is small.
	headerBarIconSize := gtk.ICON_SIZE_SMALL_TOOLBAR
	if largeToolbarThemesRegexp.MatchString(themeName) {
		headerBarIconSize = gtk.ICON_SIZE_LARGE_TOOLBAR
	}

	// Should we use symbolic icons?
	useSymbolicIcons := !nonSymbolicIconThemesRegexp.MatchString(themeName)

	// We will call icon() to automatically add -symbolic if needed.
	icon := func(s string) string {
		if useSymbolicIcons {
			return s + "-symbolic"
		}
		return s
	}

	firstImg, err := gtk.ImageNewFromIconName(icon("go-first"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.first.SetImage(firstImg)
	}

	previousImg, err := gtk.ImageNewFromIconName(icon("go-previous"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.previous.SetImage(previousImg)
	}

	nextImg, err := gtk.ImageNewFromIconName(icon("go-next"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.next.SetImage(nextImg)
	}

	newestImg, err := gtk.ImageNewFromIconName(icon("go-last"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.newest.SetImage(newestImg)
	}

	bookmarksImg, err := gtk.ImageNewFromIconName(icon("user-bookmarks"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.bookmarks.SetImage(bookmarksImg)
	}

	searchImg, err := gtk.ImageNewFromIconName(icon("edit-find"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.search.SetImage(searchImg)
	}

	menuImg, err := gtk.ImageNewFromIconName(icon("open-menu"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.menu.SetImage(menuImg)
	}
}
