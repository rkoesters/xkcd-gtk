package widget

import (
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

func NewDarkModeSwitch(setter func(darkMode bool)) (*DarkModeSwitch, error) {
	var err error

	dms := &DarkModeSwitch{}

	dms.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
	if err != nil {
		return nil, err
	}
	dms.box.SetHomogeneous(false)

	dms.swtch, err = gtk.SwitchNew()
	if err != nil {
		return nil, err
	}
	dms.swtch.SetTooltipText(l("Toggle dark mode"))
	dms.swtch.Connect("notify::active", func() {
		setter(dms.swtch.GetActive())
	})
	dms.box.PackEnd(dms.swtch, false, true, 0)

	// Use a ModelButton to force the theme to apply the same padding as the
	// other items on the WindowMenu. Has the added benefit of increasing
	// the clickable area for toggling dark mode.
	dms.label, err = gtk.ModelButtonNew()
	if err != nil {
		return nil, err
	}
	dms.label.SetLabel(l("Dark mode"))
	// Keyboard navigation should go to the switch rather than the label.
	dms.label.SetCanFocus(false)
	// Use "button-release-event" instead of "clicked" because "b-r-e" runs
	// the default handlers after our closure, whereas "clicked" always runs
	// the default handlers first. Running before the default handlers
	// allows us to prevent the running of the default handlers.
	dms.label.Connect("button-release-event", func() {
		dms.swtch.SetActive(!dms.swtch.GetActive())
		// Prevent the default handlers from running so they do not
		// close the popover menu. The popover menu should remain open
		// to mimic the behavior of clicking the switch (which would not
		// close the menu).
		dms.label.StopEmission("button-release-event")
	})
	lc, err := dms.label.GetChild()
	if err != nil {
		return nil, err
	}
	lc.ToWidget().SetHAlign(gtk.ALIGN_START)
	sc, err := dms.label.GetStyleContext()
	if err != nil {
		return nil, err
	}
	// Add this style class to disable the min-width sometimes given to
	// ModelButtons. The other ModelButtons in the menu will set the
	// min-width which gives us the flexibility to add the switch without
	// widening the menu.
	sc.AddClass(style.ClassNoMinWidth)
	dms.box.PackStart(dms.label, true, true, 0)

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

// SyncDarkMode informs the switch whether dark mode is enabled or not.
func (dms *DarkModeSwitch) SyncDarkMode(darkMode bool) {
	dms.swtch.SetActive(darkMode)
}

func (dms *DarkModeSwitch) SetCompact(compact bool) {
	if compact {
		dms.swtch.SetMarginEnd(style.PaddingPopoverCompact)
	} else {
		dms.swtch.SetMarginEnd(0)
	}
}
