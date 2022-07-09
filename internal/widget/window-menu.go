package widget

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type WindowMenu struct {
	menuButton *gtk.MenuButton

	popover    *gtk.Popover
	popoverBox *gtk.Box

	zoomBox        *ZoomBox
	darkModeSwitch *DarkModeSwitch
}

var _ Widget = &WindowMenu{}

func NewWindowMenu(prefersAppMenu bool, darkModeGetter func() bool, darkModeSetter func(bool)) (*WindowMenu, error) {
	var err error

	wm := &WindowMenu{}

	// Create the menu
	wm.menuButton, err = gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	wm.menuButton.SetTooltipText(l("Menu"))

	wm.popover, err = gtk.PopoverNew(wm.menuButton)
	if err != nil {
		return nil, err
	}
	wm.menuButton.SetPopover(wm.popover)
	wm.menuButton.SetUsePopover(true)

	wm.popoverBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	wm.popover.Add(wm.popoverBox)
	defer wm.popoverBox.ShowAll()

	addMenuSeparator := func() error {
		sep, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
		if err != nil {
			return err
		}
		wm.popoverBox.PackStart(sep, false, true, style.PaddingPopoverCompact/2)
		return nil
	}

	addMenuEntry := func(label, action string) error {
		mb, err := gtk.ModelButtonNew()
		if err != nil {
			return err
		}
		mb.SetActionName(action)
		mb.SetLabel(label)
		mbl, err := mb.GetChild()
		if err != nil {
			return err
		}
		mbl.ToWidget().SetHAlign(gtk.ALIGN_START)
		wm.popoverBox.PackStart(mb, false, true, 0)
		return nil
	}

	// Zoom section.
	wm.zoomBox, err = NewZoomBox()
	if err != nil {
		return nil, err
	}
	wm.zoomBox.box.SetMarginBottom(style.PaddingPopoverCompact / 2)
	wm.popoverBox.Add(wm.zoomBox.IWidget())

	if err = addMenuSeparator(); err != nil {
		return nil, err
	}

	// Comic properties section.
	err = addMenuEntry(l("Open link"), "win.open-link")
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(l("Explain"), "win.explain")
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(l("Properties"), "win.show-properties")
	if err != nil {
		return nil, err
	}

	// If the desktop environment will show an app menu, then we do not need
	// to add the app menu contents to the window menu.
	if prefersAppMenu {
		return wm, nil
	}

	if err = addMenuSeparator(); err != nil {
		return nil, err
	}

	err = addMenuEntry(l("New window"), "app.new-window")
	if err != nil {
		return nil, err
	}

	if err = addMenuSeparator(); err != nil {
		return nil, err
	}

	wm.darkModeSwitch, err = NewDarkModeSwitch(darkModeGetter, darkModeSetter)
	if err != nil {
		return nil, err
	}
	wm.popoverBox.PackStart(wm.darkModeSwitch.IWidget(), false, true, 0)

	if err = addMenuSeparator(); err != nil {
		return nil, err
	}

	err = addMenuEntry(l("What If?"), "app.open-what-if")
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(l("xkcd blog"), "app.open-blog")
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(l("xkcd store"), "app.open-store")
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(l("About xkcd"), "app.open-about-xkcd")
	if err != nil {
		return nil, err
	}

	if err = addMenuSeparator(); err != nil {
		return nil, err
	}

	err = addMenuEntry(l("Keyboard shortcuts"), "app.show-shortcuts")
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(l("About Comic Sticks"), "app.show-about")
	if err != nil {
		return nil, err
	}

	return wm, nil
}

func (wm *WindowMenu) Destroy() {
	wm.menuButton = nil
	wm.popover = nil
	wm.popoverBox = nil
	wm.zoomBox.Destroy()
	wm.zoomBox = nil
	wm.darkModeSwitch.Destroy()
	wm.darkModeSwitch = nil
}

func (wm *WindowMenu) IWidget() gtk.IWidget {
	return wm.menuButton
}

func (wm *WindowMenu) SetButtonImage(image gtk.IWidget) {
	wm.menuButton.SetImage(image)
}

func (wm *WindowMenu) SetCompact(compact bool) {
	if compact {
		wm.popoverBox.SetMarginTop(style.PaddingPopoverCompact)
		wm.popoverBox.SetMarginBottom(style.PaddingPopoverCompact)
		wm.popoverBox.SetMarginStart(0)
		wm.popoverBox.SetMarginEnd(0)
	} else {
		wm.popoverBox.SetMarginTop(style.PaddingPopover)
		wm.popoverBox.SetMarginBottom(style.PaddingPopover)
		wm.popoverBox.SetMarginStart(style.PaddingPopover)
		wm.popoverBox.SetMarginEnd(style.PaddingPopover)
	}
	wm.zoomBox.SetCompact(compact)
	wm.darkModeSwitch.SetCompact(compact)
}
