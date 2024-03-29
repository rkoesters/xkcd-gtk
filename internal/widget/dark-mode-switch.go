package widget

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

// DarkModeSwitch is a labeled gtk.Switch intended for toggling whether the
// application's dark mode is enabled.
type DarkModeSwitch struct {
	*gtk.Box

	label *gtk.ModelButton
	swtch *gtk.Switch

	// darkMode and setDarkMode are used to interact with the application's dark
	// mode state.
	darkMode    func() bool
	setDarkMode func(bool)
}

var _ Widget = &DarkModeSwitch{}

func NewDarkModeSwitch(darkModeGetter func() bool, darkModeSetter func(bool)) (*DarkModeSwitch, error) {
	super, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 2)
	if err != nil {
		return nil, err
	}
	dms := &DarkModeSwitch{
		Box: super,

		darkMode:    darkModeGetter,
		setDarkMode: darkModeSetter,
	}

	dms.SetHomogeneous(false)

	dms.swtch, err = gtk.SwitchNew()
	if err != nil {
		return nil, err
	}
	dms.swtch.SetTooltipText(l("Toggle dark mode"))
	dms.swtch.SetActive(dms.darkMode())
	dms.swtch.Connect("notify::active", dms.SwitchStateChanged)
	dms.PackEnd(dms.swtch, false, true, 0)

	// Use a ModelButton to force the theme to apply the same padding as the
	// other items on the WindowMenu. Has the added benefit of increasing the
	// clickable area for toggling dark mode.
	dms.label, err = gtk.ModelButtonNew()
	if err != nil {
		return nil, err
	}
	dms.label.SetLabel(l("Dark mode"))
	// Keyboard navigation should go to the switch rather than the label.
	dms.label.SetCanFocus(false)
	// Use "button-release-event" instead of "clicked" because "b-r-e" runs the
	// default handlers after our closure, whereas "clicked" always runs the
	// default handlers first. Running before the default handlers allows us to
	// prevent the running of the default handlers.
	dms.label.Connect("button-release-event", func(w gtk.IWidget) {
		dms.swtch.Activate()
		// Prevent the default handlers from running so they do not close the
		// popover menu. The popover menu should remain open to mimic the
		// behavior of clicking the switch (which would not close the menu).
		w.ToWidget().StopEmission("button-release-event")
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
	// ModelButtons. The other ModelButtons in the menu will set the min-width
	// which gives us the flexibility to add the switch without widening the
	// menu.
	sc.AddClass(style.ClassNoMinWidth)
	dms.PackStart(dms.label, true, true, 0)

	dms.ShowAll()
	return dms, nil
}

func (dms *DarkModeSwitch) Dispose() {
	if dms == nil {
		return
	}

	dms.Box = nil

	dms.label = nil
	dms.swtch = nil
	dms.darkMode = nil
	dms.setDarkMode = nil
}

// SwitchStateChanged is called when the active state of the switch changes.
func (dms *DarkModeSwitch) SwitchStateChanged(swtch *gtk.Switch) {
	swtchState := swtch.GetActive()
	// Avoid calling dms.setDarkMode when this signal might have been emitted by
	// dms.SyncDarkMode.
	if swtchState == dms.darkMode() {
		return
	}
	dms.setDarkMode(swtchState)
}

// SyncDarkMode informs the switch whether dark mode is enabled or not.
func (dms *DarkModeSwitch) SyncDarkMode(darkMode bool) {
	if dms == nil {
		return
	}
	// Avoid calling dms.swtch.SetActive when this signal might have been
	// emitted by dms.SwitchStateChanged.
	if darkMode == dms.swtch.GetActive() {
		return
	}
	dms.swtch.SetActive(darkMode)
}

func (dms *DarkModeSwitch) SetCompact(compact bool) {
	if dms == nil {
		return
	}
	if compact {
		dms.swtch.SetMarginEnd(style.PaddingPopoverCompact)
	} else {
		dms.swtch.SetMarginEnd(0)
	}
}
