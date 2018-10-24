package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
)

const css = `
@define-color colorPrimary #96a8c8;
@define-color textColorPrimary #1a1a1a;

.comic-container > .frame {
	background-color: #ffffff;
}
`

var (
	// largeToolbarThemes is the list of gtk themes for which we should use
	// large toolbar buttons.
	largeToolbarThemes = []string{
		"elementary",
		"win32",
	}

	// nonSymbolicIconThemes is the list of gtk themes for which we
	// should use non-symbolic icons.
	nonSymbolicIconThemes = []string{
		"elementary",
	}

	// skinnyMenuThemes is the list of gtk themes for which we should
	// use skinny popover menus.
	skinnyMenuThemes = []string{
		"elementary",
	}
)

// LoadCSS provides the application's custom CSS to GTK.
func (app *Application) LoadCSS() {
	provider, err := gtk.CssProviderNew()
	if err != nil {
		log.Print(err)
		return
	}
	provider.LoadFromData(css)

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
		settings, err := gtk.SettingsGetDefault()
		if err != nil {
			log.Print(err)
		} else {
			// settings.GetProperty returns an interface{}, we will convert
			// it to a string in a moment.
			themeNameIface, err := settings.GetProperty("gtk-theme-name")
			if err != nil {
				log.Print(err)
			} else {
				themeNameStr, ok := themeNameIface.(string)
				if ok {
					themeName = themeNameStr
				}
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

	nextImg, err := gtk.ImageNewFromIconName(icon("go-next"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.next.SetImage(nextImg)
	}

	previousImg, err := gtk.ImageNewFromIconName(icon("go-previous"), headerBarIconSize)
	if err != nil {
		log.Print(err)
	} else {
		win.previous.SetImage(previousImg)
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

	menuPopoverChild, err := win.menu.GetPopover().GetChild()
	if err != nil {
		log.Print(err)
	} else {
		menuBox := (&gtk.Stack{gtk.Container{*menuPopoverChild}}).GetVisibleChild()
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
