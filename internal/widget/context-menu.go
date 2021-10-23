package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type ContextMenu struct {
	menu *gtk.Menu

	imageViewer *ImageViewer
}

var _ Widget = &ContextMenu{}

func NewContextMenu(window *gtk.ApplicationWindow, iv *ImageViewer) (*ContextMenu, error) {
	var err error

	cm := &ContextMenu{
		imageViewer: iv,
	}

	menuModel := glib.MenuNew()

	menuModel.AppendSectionWithoutLabel(&NewContextMenuSection().MenuModel)

	cm.menu, err = gtk.GtkMenuNewFromModel(&menuModel.MenuModel)
	if err != nil {
		return nil, err
	}
	cm.menu.SetHAlign(gtk.ALIGN_START)
	cm.menu.ShowAll()

	cm.menu.InsertActionGroup("win", window.IActionGroup)
	cm.menu.HideOnDelete()

	return cm, nil
}

func NewContextMenuSection() *glib.Menu {
	contextSection := glib.MenuNew()
	contextSection.Append(l("Open Link"), "win.open-link")
	contextSection.Append(l("Explain"), "win.explain")
	contextSection.Append(l("Properties"), "win.show-properties")
	return contextSection
}

func (cm *ContextMenu) Present(event *gdk.Event) {
	cm.menu.PopupAtPointer(event)
}

func (cm *ContextMenu) Destroy() {
	cm.imageViewer = nil
	cm.menu = nil
}

func (cm *ContextMenu) IWidget() gtk.IWidget {
	return cm.menu
}
