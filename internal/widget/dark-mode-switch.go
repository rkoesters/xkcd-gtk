package widget

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

// DarkModeSwitch is a labeled gtk.Switch intended for toggling whether the
// application's dark mode is enabled.
type DarkModeSwitch struct {
	box   *gtk.Box
	label *gtk.Label
	swtch *gtk.Switch
}

var _ Widget = &DarkModeSwitch{}

func NewDarkModeSwitch(setter func(active bool)) (*DarkModeSwitch, error) {
	var err error

	dms := &DarkModeSwitch{}

	dms.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, style.PopoverPaddingCompact/2)
	if err != nil {
		return nil, err
	}
	dms.box.SetHomogeneous(false)
	sc, err := dms.box.GetStyleContext()
	if err != nil {
		return nil, err
	}
	// These 3 style classes are needed to make Adwaita and elementary's
	// stylesheets use the same padding for this gtk.Box as they do for the
	// gtk.ModelButtons that WindowMenu uses.
	sc.AddClass("menuitem")
	sc.AddClass("button")
	sc.AddClass("flat")

	dms.label, err = gtk.LabelNew(l("Dark mode"))
	if err != nil {
		return nil, err
	}
	dms.label.SetHAlign(gtk.ALIGN_START)
	dms.box.PackStart(dms.label, false, true, 0)

	dms.swtch, err = gtk.SwitchNew()
	if err != nil {
		return nil, err
	}
	dms.swtch.SetTooltipText(l("Toggle dark mode"))
	dms.swtch.Connect("notify::active", func() {
		active := dms.swtch.GetActive()
		go glib.IdleAdd(func() { setter(active) })
	})
	dms.box.PackEnd(dms.swtch, false, true, 0)

	dms.box.ShowAll()
	return dms, nil
}

func (dms *DarkModeSwitch) IWidget() gtk.IWidget {
	return dms.box
}

func (dms *DarkModeSwitch) Destroy() {
	dms.box = nil
	dms.label = nil
	dms.swtch = nil
}

// SetActive informs the switch whether dark mode is enabled or not.
func (dms *DarkModeSwitch) SetActive(active bool) {
	dms.swtch.SetActive(active)
}
