package widget

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type WindowMenu struct {
	menuButton *gtk.MenuButton

	showProperties func() // win.ShowProperties
}

var _ Widget = &WindowMenu{}

func NewWindowMenu(prefersAppMenu bool) (*WindowMenu, error) {
	var err error

	wm := &WindowMenu{}

	// Create the menu
	wm.menuButton, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	wm.menuButton.SetTooltipText(l("Menu"))

	menu := glib.MenuNew()

	menu.AppendSectionWithoutLabel(&NewContextMenuSection().MenuModel)

	if !prefersAppMenu {
		appSection := glib.MenuNew()
		appSection.Append(l("New window"), "app.new-window")
		appSection.Append(l("Toggle dark mode"), "app.toggle-dark-mode")
		menu.AppendSectionWithoutLabel(&appSection.MenuModel)

		websiteSection := glib.MenuNew()
		websiteSection.Append(l("What If?"), "app.open-what-if")
		websiteSection.Append(l("XKCD Blog"), "app.open-blog")
		websiteSection.Append(l("XKCD Store"), "app.open-store")
		websiteSection.Append(l("About XKCD"), "app.open-about-xkcd")
		menu.AppendSectionWithoutLabel(&websiteSection.MenuModel)

		helpSection := glib.MenuNew()
		helpSection.Append(l("Keyboard shortcuts"), "app.show-shortcuts")
		helpSection.Append(l("About Comic Sticks"), "app.show-about")
		menu.AppendSectionWithoutLabel(&helpSection.MenuModel)
	}

	menuWidget, err := gtk.GtkMenuNewFromModel(&menu.MenuModel)
	if err != nil {
		return nil, err
	}
	menuWidget.SetHAlign(gtk.ALIGN_END)
	wm.menuButton.SetPopup(menuWidget)

	return wm, nil
}

func (wm *WindowMenu) Destroy() {
	wm.menuButton = nil
}

func (wm *WindowMenu) IWidget() gtk.IWidget {
	return wm.menuButton
}
