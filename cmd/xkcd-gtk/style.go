package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
	"log"
	"os"
)

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
	if style.IsLargeToolbarTheme(themeName) {
		headerBarIconSize = gtk.ICON_SIZE_LARGE_TOOLBAR
	}

	useSymbolicIcons := style.IsSymbolicIconTheme(themeName)

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

	randomImg, err := gtk.ImageNewFromIconName(icon("media-playlist-shuffle"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.random.SetImage(randomImg)
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
		win.bookmarksMenu.IWidget().(*gtk.MenuButton).SetImage(bookmarksImg)
	}

	searchImg, err := gtk.ImageNewFromIconName(icon("edit-find"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.searchMenu.IWidget().(*gtk.MenuButton).SetImage(searchImg)
	}

	menuImg, err := gtk.ImageNewFromIconName(icon("open-menu"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.menu.SetImage(menuImg)
	}
}
