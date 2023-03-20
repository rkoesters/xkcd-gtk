package widget

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type PopoverMenu struct {
	*gtk.Popover

	box *gtk.Box
}

var _ Widget = &PopoverMenu{}

func NewPopoverMenu(relative gtk.IWidget) (*PopoverMenu, error) {
	super, err := gtk.PopoverNew(relative)
	if err != nil {
		return nil, err
	}
	pm := &PopoverMenu{
		Popover: super,
	}

	pm.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return nil, err
	}
	pm.Add(pm.box)

	return pm, nil
}

func (pm *PopoverMenu) Dispose() {
	pm.Popover = nil
	pm.box = nil
}

func (pm *PopoverMenu) AddSeparator() error {
	sep, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return err
	}
	pm.AddChild(sep, style.PaddingPopoverCompact/2)
	return nil
}

func (pm *PopoverMenu) AddMenuEntry(label, action string) (*gtk.ModelButton, error) {
	mb, err := gtk.ModelButtonNew()
	if err != nil {
		return nil, err
	}
	mb.SetActionName(action)
	mb.SetLabel(label)
	mbl, err := mb.GetChild()
	if err != nil {
		return nil, err
	}
	mbl.ToWidget().SetHAlign(gtk.ALIGN_START)
	pm.AddChild(mb, 0)
	return mb, nil
}

func (pm *PopoverMenu) AddMenuEntries(labelActionPairs [][2]string) error {
	var err error
	for _, pair := range labelActionPairs {
		label, action := pair[0], pair[1]
		if action == "sep" {
			err = pm.AddSeparator()
			if err != nil {
				return err
			}
			continue
		}

		_, err = pm.AddMenuEntry(label, action)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pm *PopoverMenu) AddCheckButton(label string, stateGetter func() bool, stateSetter func(bool)) (*CheckModelButton, error) {
	mb, err := NewCheckModelButton(stateGetter, stateSetter)
	if err != nil {
		return nil, err
	}
	mb.SetLabel(label)
	mbl, err := mb.GetChild()
	if err != nil {
		return nil, err
	}
	mbl.ToWidget().SetHAlign(gtk.ALIGN_START)
	pm.AddChild(mb, 0)
	return mb, nil
}

func (pm *PopoverMenu) AddChild(child gtk.IWidget, padding uint) {
	pm.box.PackStart(child, false, true, padding)
}

func (pm *PopoverMenu) ShowAll() {
	pm.box.ShowAll()
}

func (pm *PopoverMenu) SetCompact(compact bool) {
	if compact {
		pm.box.SetMarginTop(style.PaddingPopoverCompact)
		pm.box.SetMarginBottom(style.PaddingPopoverCompact)
		pm.box.SetMarginStart(0)
		pm.box.SetMarginEnd(0)
	} else {
		pm.box.SetMarginTop(style.PaddingPopover)
		pm.box.SetMarginBottom(style.PaddingPopover)
		pm.box.SetMarginStart(style.PaddingPopover)
		pm.box.SetMarginEnd(style.PaddingPopover)
	}
}
