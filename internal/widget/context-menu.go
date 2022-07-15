package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type ContextMenu struct {
	*gtk.Menu
}

var _ Widget = &ContextMenu{}

func NewContextMenu(actionGroup glib.IActionGroup) (*ContextMenu, error) {
	menuModel := glib.MenuNew()

	bookmarkSection := glib.MenuNew()
	bookmarkSection.Append(l("Add to bookmarks"), "win.bookmark-new")
	bookmarkSection.Append(l("Remove from bookmarks"), "win.bookmark-remove")
	menuModel.AppendSectionWithoutLabel(&bookmarkSection.MenuModel)

	zoomSection := glib.MenuNew()
	zoomSection.Append(l("Zoom in"), "win.zoom-in")
	zoomSection.Append(l("Zoom out"), "win.zoom-out")
	zoomSection.Append(l("Reset zoom"), "win.zoom-reset")
	menuModel.AppendSectionWithoutLabel(&zoomSection.MenuModel)

	propertiesSection := glib.MenuNew()
	propertiesSection.Append(l("Open link"), "win.open-link")
	propertiesSection.Append(l("Explain"), "win.explain")
	propertiesSection.Append(l("Properties"), "win.show-properties")
	menuModel.AppendSectionWithoutLabel(&propertiesSection.MenuModel)

	super, err := gtk.GtkMenuNewFromModel(&menuModel.MenuModel)
	if err != nil {
		return nil, err
	}
	cm := &ContextMenu{
		Menu: super,
	}

	cm.SetHAlign(gtk.ALIGN_START)
	cm.ShowAll()

	cm.InsertActionGroup("win", actionGroup)
	cm.HideOnDelete()

	return cm, nil
}

func (cm *ContextMenu) Present(event *gdk.Event) {
	cm.PopupAtPointer(event)
}

func (cm *ContextMenu) Dispose() {
	if cm == nil {
		return
	}

	cm.Menu = nil
}
