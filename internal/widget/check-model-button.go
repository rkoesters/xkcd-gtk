package widget

import (
	"github.com/gotk3/gotk3/gtk"
)

type CheckModelButton struct {
	*gtk.ModelButton

	// state and setState are used to interact with the source of truth for the
	// state of this boolean.
	state    func() bool
	setState func(bool)
}

var _ Widget = &CheckModelButton{}

func NewCheckModelButton(stateGetter func() bool, stateSetter func(bool)) (*CheckModelButton, error) {
	super, err := gtk.ModelButtonNew()
	if err != nil {
		return nil, err
	}
	cmb := &CheckModelButton{
		ModelButton: super,

		state:    stateGetter,
		setState: stateSetter,
	}
	cmb.SetProperty("role", gtk.BUTTON_ROLE_CHECK)
	cmb.Connect("clicked", cmb.Clicked)

	cmb.SyncState(cmb.state())

	return cmb, nil
}

func (cmb *CheckModelButton) Dispose() {
	cmb.ModelButton = nil
}

func (cmb *CheckModelButton) Clicked() {
	cmb.setState(!cmb.state())
}

func (cmb *CheckModelButton) SyncState(state bool) {
	cmb.SetProperty("active", state)
}
