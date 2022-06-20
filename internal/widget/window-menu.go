package widget

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type WindowMenu struct {
	menuButton *gtk.MenuButton

	popover    *gtk.Popover
	popoverBox *gtk.Box

	zoomBox *ZoomBox

	showProperties func() // win.ShowProperties
}

var _ Widget = &WindowMenu{}

func NewWindowMenu(accels *gtk.AccelGroup, comicContainer *ImageViewer, prefersAppMenu bool) (*WindowMenu, error) {
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
	wm.popoverBox.SetMarginTop(style.PopoverMenuPadding)
	wm.popoverBox.SetMarginBottom(style.PopoverMenuPadding)
	wm.popoverBox.SetMarginStart(style.PopoverMenuPadding)
	wm.popoverBox.SetMarginEnd(style.PopoverMenuPadding)
	wm.popover.Add(wm.popoverBox)

	addMenuSeparator := func(menuBox *gtk.Box) error {
		sep, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
		if err != nil {
			return err
		}
		menuBox.PackStart(sep, false, true, style.PopoverMenuPadding/2)
		return nil
	}

	addMenuEntry := func(menuBox *gtk.Box, label, action string) error {
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
		menuBox.PackStart(mb, false, true, 0)
		return nil
	}

	// Zoom section.
	wm.zoomBox, err = NewZoomBox(accels, comicContainer)
	if err != nil {
		return nil, err
	}
	wm.zoomBox.SetCurrentZoom(comicContainer.scale)
	wm.popoverBox.Add(wm.zoomBox.IWidget())
	err = addMenuSeparator(wm.popoverBox)
	if err != nil {
		return nil, err
	}

	// Comic properties section.
	err = addMenuEntry(wm.popoverBox, l("Open link"), "win.open-link")
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(wm.popoverBox, l("Explain"), "win.explain")
	if err != nil {
		return nil, err
	}
	err = addMenuEntry(wm.popoverBox, l("Properties"), "win.show-properties")
	if err != nil {
		return nil, err
	}
	if err = addMenuSeparator(wm.popoverBox); err != nil {
		return nil, err
	}

	if !prefersAppMenu {
		err = addMenuEntry(wm.popoverBox, l("New window"), "app.new-window")
		if err != nil {
			return nil, err
		}
		err = addMenuEntry(wm.popoverBox, l("Toggle dark mode"), "app.toggle-dark-mode")
		if err != nil {
			return nil, err
		}
		if err = addMenuSeparator(wm.popoverBox); err != nil {
			return nil, err
		}

		err = addMenuEntry(wm.popoverBox, l("What If?"), "app.open-what-if")
		if err != nil {
			return nil, err
		}
		err = addMenuEntry(wm.popoverBox, l("XKCD Blog"), "app.open-blog")
		if err != nil {
			return nil, err
		}
		err = addMenuEntry(wm.popoverBox, l("XKCD Store"), "app.open-store")
		if err != nil {
			return nil, err
		}
		err = addMenuEntry(wm.popoverBox, l("About XKCD"), "app.open-about-xkcd")
		if err != nil {
			return nil, err
		}
		if err = addMenuSeparator(wm.popoverBox); err != nil {
			return nil, err
		}

		err = addMenuEntry(wm.popoverBox, l("Keyboard shortcuts"), "app.show-shortcuts")
		if err != nil {
			return nil, err
		}
		err = addMenuEntry(wm.popoverBox, l("About Comic Sticks"), "app.show-about")
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
}

func (wm *WindowMenu) IWidget() gtk.IWidget {
	return wm.menuButton
}
