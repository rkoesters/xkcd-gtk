package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type WindowMenu struct {
	*gtk.MenuButton

	popover *PopoverMenu

	zoomBox        *ZoomBox
	darkModeSwitch *DarkModeSwitch // may be nil
}

var _ Widget = &WindowMenu{}

func NewWindowMenu(accels *gtk.AccelGroup, prefersAppMenu bool, darkModeGetter func() bool, darkModeSetter func(bool)) (*WindowMenu, error) {
	super, err := gtk.MenuButtonNew()
	if err != nil {
		return nil, err
	}
	wm := &WindowMenu{
		MenuButton: super,
	}

	wm.SetTooltipText(l("Window menu"))
	wm.AddAccelerator("activate", accels, gdk.KEY_F10, 0, gtk.ACCEL_VISIBLE)

	wm.popover, err = NewPopoverMenu(wm)
	if err != nil {
		return nil, err
	}
	wm.popover.SetSizeRequest(200, -1)
	wm.SetPopover(wm.popover.Popover)
	wm.SetUsePopover(true)

	defer wm.popover.ShowAll()

	// Zoom section.
	wm.zoomBox, err = NewZoomBox()
	if err != nil {
		return nil, err
	}
	wm.zoomBox.SetMarginBottom(style.PaddingPopoverCompact / 2)
	wm.popover.AddChild(wm.zoomBox, style.PaddingPopoverCompact/2)

	err = wm.popover.AddMenuEntries([][2]string{
		{"", "sep"},
		{l("Open link"), "win.open-link"},
		{l("Explain"), "win.explain"},
		{l("Properties"), "win.show-properties"},
	})
	if err != nil {
		return nil, err
	}

	// If the desktop environment will show an app menu, then we do not need to
	// add the app menu contents to the window menu.
	if prefersAppMenu {
		return wm, nil
	}

	err = wm.popover.AddMenuEntries([][2]string{
		{"", "sep"},
		{l("New window"), "app.new-window"},
		{"", "sep"},
		{l("Cache manager"), "app.show-cache"},
		{"", "sep"},
	})
	if err != nil {
		return nil, err
	}

	wm.darkModeSwitch, err = NewDarkModeSwitch(darkModeGetter, darkModeSetter)
	if err != nil {
		return nil, err
	}
	wm.popover.AddChild(wm.darkModeSwitch, 0)

	err = wm.popover.AddMenuEntries([][2]string{
		{"", "sep"},
		{l("What If?"), "app.open-what-if"},
		{l("xkcd blog"), "app.open-blog"},
		{l("xkcd books"), "app.open-books"},
		{l("About xkcd"), "app.open-about-xkcd"},
		{"", "sep"},
		{l("Keyboard shortcuts"), "app.show-shortcuts"},
		{l("About"), "app.show-about"},
	})
	if err != nil {
		return nil, err
	}

	return wm, nil
}

func (wm *WindowMenu) Dispose() {
	if wm == nil {
		return
	}

	wm.MenuButton = nil

	wm.popover.Dispose()
	wm.popover = nil
	wm.zoomBox.Dispose()
	wm.zoomBox = nil
	wm.darkModeSwitch.Dispose()
	wm.darkModeSwitch = nil
}

func (wm *WindowMenu) SetCompact(compact bool) {
	wm.popover.SetCompact(compact)
	wm.zoomBox.SetCompact(compact)
	wm.darkModeSwitch.SetCompact(compact)
}
