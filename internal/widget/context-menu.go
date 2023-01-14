package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type ContextMenu struct {
	*PopoverMenu

	zoomBox *ZoomBox
}

var _ Widget = &ContextMenu{}

func NewContextMenu(relative gtk.IWidget, actionGroup glib.IActionGroup) (*ContextMenu, error) {
	super, err := NewPopoverMenu(relative)
	if err != nil {
		return nil, err
	}
	cm := &ContextMenu{
		PopoverMenu: super,
	}

	defer cm.ShowAll()

	err = cm.AddMenuEntry(l("Add to bookmarks"), "win.bookmark-new")
	if err != nil {
		return nil, err
	}
	err = cm.AddMenuEntry(l("Remove from bookmarks"), "win.bookmark-remove")
	if err != nil {
		return nil, err
	}

	if err = cm.AddSeparator(); err != nil {
		return nil, err
	}

	// zoom
	cm.zoomBox, err = NewZoomBox()
	if err != nil {
		return nil, err
	}
	cm.zoomBox.SetMarginBottom(style.PaddingPopoverCompact / 2)
	cm.zoomBox.SetMarginTop(style.PaddingPopoverCompact / 2)
	//cm.AddChild(cm.zoomBox, style.PaddingPopoverCompact/2)
	//cm.zoomBox.SetMarginBottom(0)
	cm.AddChild(cm.zoomBox, 0)

	if err = cm.AddSeparator(); err != nil {
		return nil, err
	}

	err = cm.AddMenuEntry(l("Open link"), "win.open-link")
	if err != nil {
		return nil, err
	}
	err = cm.AddMenuEntry(l("Explain"), "win.explain")
	if err != nil {
		return nil, err
	}
	err = cm.AddMenuEntry(l("Properties"), "win.show-properties")
	if err != nil {
		return nil, err
	}

	cm.InsertActionGroup("win", actionGroup)
	cm.HideOnDelete()
	cm.SetModal(true)
	cm.SetPosition(gtk.POS_BOTTOM)

	return cm, nil
}

func (cm *ContextMenu) Dispose() {
	cm.PopoverMenu.Dispose()
	cm.PopoverMenu = nil
	cm.zoomBox.Dispose()
	cm.zoomBox = nil
}

func (cm *ContextMenu) PopupAtPointer(event *gdk.EventButton) {
	rect := gdk.RectangleNew(int(event.X()), int(event.Y()), 0, 0)
	cm.SetPointingTo(*rect)
	cm.Popup()
}

func (cm *ContextMenu) SetCompact(compact bool) {
	cm.PopoverMenu.SetCompact(compact)
	cm.zoomBox.SetCompact(compact)
}
