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

	showProperties func() // win.ShowProperties
}

var _ Widget = &WindowMenu{}

func NewWindowMenu(prefersAppMenu bool, setDarkMode func(bool)) (*WindowMenu, error) {
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
	wm.popoverBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	wm.popover.Add(wm.popoverBox)

	addMenuSeparator := func() error {
		sep, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
		if err != nil {
			return err
		}
		wm.popoverBox.PackStart(sep, false, true, style.PopoverPaddingCompact/2)
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
	wm.zoomBox.box.SetMarginBottom(style.PopoverPaddingCompact / 2)
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

	if !prefersAppMenu {
		if err = addMenuSeparator(); err != nil {
			return nil, err
		}

		err = addMenuEntry(l("New window"), "app.new-window")
		if err != nil {
			return nil, err
		}
		wm.darkModeSwitch, err = NewDarkModeSwitch(setDarkMode)
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
		err = addMenuEntry(l("XKCD Blog"), "app.open-blog")
		if err != nil {
			return nil, err
		}
		err = addMenuEntry(l("XKCD Store"), "app.open-store")
		if err != nil {
			return nil, err
		}
		err = addMenuEntry(l("About XKCD"), "app.open-about-xkcd")
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
	}

	wm.popoverBox.ShowAll()
	wm.menuButton.SetPopover(wm.popover)
	wm.menuButton.SetUsePopover(true)

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
		wm.popoverBox.SetMarginTop(style.PopoverPaddingCompact)
		wm.popoverBox.SetMarginBottom(style.PopoverPaddingCompact)
		wm.popoverBox.SetMarginStart(0)
		wm.popoverBox.SetMarginEnd(0)
	} else {
		wm.popoverBox.SetMarginTop(style.PopoverPadding)
		wm.popoverBox.SetMarginBottom(style.PopoverPadding)
		wm.popoverBox.SetMarginStart(style.PopoverPadding)
		wm.popoverBox.SetMarginEnd(style.PopoverPadding)
	}
}
