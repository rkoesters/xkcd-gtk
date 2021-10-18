package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type WindowMenu struct {
	actions map[string]*glib.SimpleAction // ptr to win.actions
	accels  *gtk.AccelGroup               // ptr to win.accels

	menuButton *gtk.MenuButton

	showProperties func() // win.ShowProperties
}

var _ Widget = &WindowMenu{}

func NewWindowMenu(prefersAppMenu bool, actions map[string]*glib.SimpleAction, accels *gtk.AccelGroup, propertiesShower func()) (*WindowMenu, error) {
	var err error

	wm := &WindowMenu{
		actions:        actions,
		accels:         accels,
		showProperties: propertiesShower,
	}

	// Create the menu
	wm.menuButton, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	wm.menuButton.SetTooltipText(l("Menu"))

	menu := glib.MenuNew()

	menuSection1 := glib.MenuNew()
	menuSection1.Append(l("Open Link"), "win.open-link")
	menuSection1.Append(l("Explain"), "win.explain")
	menuSection1.Append(l("Properties"), "win.show-properties")
	menu.AppendSectionWithoutLabel(&menuSection1.MenuModel)
	wm.accels.Connect(gdk.KEY_p, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE, wm.showProperties)

	if !prefersAppMenu {
		menuSection2 := glib.MenuNew()
		menuSection2.Append(l("New Window"), "app.new-window")
		menu.AppendSectionWithoutLabel(&menuSection2.MenuModel)

		menuSection3 := glib.MenuNew()
		menuSection3.Append(l("Toggle Dark Mode"), "app.toggle-dark-mode")
		menu.AppendSectionWithoutLabel(&menuSection3.MenuModel)

		menuSection4 := glib.MenuNew()
		menuSection4.Append(l("What If?"), "app.open-what-if")
		menuSection4.Append(l("XKCD Blog"), "app.open-blog")
		menuSection4.Append(l("XKCD Store"), "app.open-store")
		menuSection4.Append(l("About XKCD"), "app.open-about-xkcd")
		menu.AppendSectionWithoutLabel(&menuSection4.MenuModel)

		menuSection5 := glib.MenuNew()
		menuSection5.Append(l("Keyboard Shortcuts"), "app.show-shortcuts")
		menuSection5.Append(l("About Comic Sticks"), "app.show-about")
		menu.AppendSectionWithoutLabel(&menuSection5.MenuModel)
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
	wm.actions = nil
	wm.accels = nil

	wm.menuButton = nil
}

func (wm *WindowMenu) IWidget() gtk.IWidget {
	return wm.menuButton
}
