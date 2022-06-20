package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type ContextMenu struct {
	menu *gtk.Menu
}

var _ Widget = &ContextMenu{}

func NewContextMenu(actionGroup glib.IActionGroup) (*ContextMenu, error) {
	var err error

	cm := &ContextMenu{}

	menuModel := glib.MenuNew()

	menuModel.AppendSectionWithoutLabel(&NewBookmarkMenuSection().MenuModel)

	menuModel.AppendSectionWithoutLabel(&NewZoomMenuSection().MenuModel)

	menuModel.AppendSectionWithoutLabel(&NewComicPropertiesMenuSection().MenuModel)

	cm.menu, err = gtk.GtkMenuNewFromModel(&menuModel.MenuModel)
	if err != nil {
		return nil, err
	}
	cm.menu.SetHAlign(gtk.ALIGN_START)
	cm.menu.ShowAll()

	cm.menu.InsertActionGroup("win", actionGroup)
	cm.menu.HideOnDelete()

	return cm, nil
}

func NewBookmarkMenuSection() *glib.Menu {
	bookmarkSection := glib.MenuNew()
	bookmarkSection.Append(l("Add to bookmarks"), "win.bookmark-new")
	bookmarkSection.Append(l("Remove from bookmarks"), "win.bookmark-remove")
	return bookmarkSection
}

func NewZoomMenuSection() *glib.Menu {
	zoomSection := glib.MenuNew()
	zoomSection.Append(l("Zoom in"), "win.zoom-in")
	zoomSection.Append(l("Zoom out"), "win.zoom-out")
	zoomSection.Append(l("Reset zoom"), "win.zoom-reset")
	return zoomSection
}

func NewComicPropertiesMenuSection() *glib.Menu {
	propertiesSection := glib.MenuNew()
	propertiesSection.Append(l("Open link"), "win.open-link")
	propertiesSection.Append(l("Explain"), "win.explain")
	propertiesSection.Append(l("Properties"), "win.show-properties")
	return propertiesSection
}

func (cm *ContextMenu) Present(event *gdk.Event) {
	cm.menu.PopupAtPointer(event)
}

func (cm *ContextMenu) Destroy() {
	cm.menu = nil
}

func (cm *ContextMenu) IWidget() gtk.IWidget {
	return cm.menu
}
