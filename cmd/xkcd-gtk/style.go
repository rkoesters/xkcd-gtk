package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const (
	styleClassComicContainer = "comic-container"
	styleClassDark           = "dark"
	styleClassLinked         = "linked"
)

var (
	// largeToolbarThemes is the list of gtk themes for which we should use
	// large toolbar buttons.
	largeToolbarThemes = []string{
		"elementary",
		"elementary-x",
		"win32",
	}

	// nonSymbolicIconThemes is the list of gtk themes for which we should
	// use non-symbolic icons.
	nonSymbolicIconThemes = []string{
		"elementary",
		"elementary-x",
	}

	// skinnyMenuThemes is the list of gtk themes for which we should use
	// skinny popover menus.
	skinnyMenuThemes = []string{
		"elementary",
		"elementary-x",
	}
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
	for _, largeToolbarTheme := range largeToolbarThemes {
		if themeName == largeToolbarTheme {
			headerBarIconSize = gtk.ICON_SIZE_LARGE_TOOLBAR
			break
		}
	}

	// Should we use symbolic icons?
	useSymbolicIcons := true
	for _, nonSymbolicIconTheme := range nonSymbolicIconThemes {
		if themeName == nonSymbolicIconTheme {
			useSymbolicIcons = false
			break
		}
	}

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

	// Should we use skinny popover menus?
	useSkinnyMenus := false
	for _, skinnyMenuTheme := range skinnyMenuThemes {
		if themeName == skinnyMenuTheme {
			useSkinnyMenus = true
			break
		}
	}

	c, err := win.menu.GetPopover().GetChild()
	if err != nil {
		log.Print(err)
	} else {
		menuBox := c.ToWidget()
		if useSkinnyMenus {
			menuBox.SetMarginTop(4)
			menuBox.SetMarginBottom(4)
			menuBox.SetMarginStart(0)
			menuBox.SetMarginEnd(0)
		} else {
			menuBox.SetMarginTop(10)
			menuBox.SetMarginBottom(10)
			menuBox.SetMarginStart(10)
			menuBox.SetMarginEnd(10)
		}
	}
}
