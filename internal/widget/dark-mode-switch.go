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
	label *gtk.ModelButton
	swtch *gtk.Switch
}

var _ Widget = &DarkModeSwitch{}

func NewDarkModeSwitch(setter func(active bool)) (*DarkModeSwitch, error) {
	var err error

	dms := &DarkModeSwitch{}

	dms.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, style.PopoverPaddingCompact)
	if err != nil {
		return nil, err
	}
	dms.box.SetHomogeneous(false)

	// Use a ModelButton to force the theme to apply the same padding as the
	// other items on the WindowMenu. Has the added benefit of increasing
	// the clickable area for toggling dark mode.
	dms.label, err = gtk.ModelButtonNew()
	if err != nil {
		return nil, err
	}
	dms.label.SetLabel(l("Dark mode"))
	dms.label.SetActionName("app.toggle-dark-mode")
	dms.label.SetTooltipText(l("Toggle dark mode"))
	lc, err := dms.label.GetChild()
	if err != nil {
		return nil, err
	}
	lc.ToWidget().SetHAlign(gtk.ALIGN_START)
	sc, err := dms.label.GetStyleContext()
	if err != nil {
		return nil, err
	}
	sc.AddClass(style.ClassNarrowModelButton)
	dms.box.PackStart(dms.label, true, true, 0)

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
