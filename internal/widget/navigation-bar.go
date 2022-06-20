package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type NavigationBar struct {
	accels *gtk.AccelGroup // ptr to win.accels

	box *gtk.ButtonBox

	firstButton    *gtk.Button
	previousButton *gtk.Button
	randomButton   *gtk.Button
	nextButton     *gtk.Button
	newestButton   *gtk.Button
}

var _ Widget = &NavigationBar{}

func NewNavigationBar(accels *gtk.AccelGroup) (*NavigationBar, error) {
	var err error

	nb := &NavigationBar{
		accels: accels,
	}

	nb.box, err = gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	nb.box.SetLayout(gtk.BUTTONBOX_EXPAND)

	nb.firstButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	nb.firstButton.SetTooltipText(l("Go to the first comic"))
	nb.firstButton.SetProperty("action-name", "win.first-comic")
	nb.firstButton.AddAccelerator("activate", nb.accels, gdk.KEY_Home, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
	nb.box.Add(nb.firstButton)

	nb.previousButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	nb.previousButton.SetTooltipText(l("Go to the previous comic"))
	nb.previousButton.SetProperty("action-name", "win.previous-comic")
	nb.previousButton.AddAccelerator("activate", nb.accels, gdk.KEY_Left, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
	nb.box.Add(nb.previousButton)

	nb.randomButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	nb.randomButton.SetTooltipText(l("Go to a random comic"))
	nb.randomButton.SetProperty("action-name", "win.random-comic")
	nb.randomButton.AddAccelerator("activate", nb.accels, gdk.KEY_r, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
	nb.box.Add(nb.randomButton)

	nb.nextButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	nb.nextButton.SetTooltipText(l("Go to the next comic"))
	nb.nextButton.SetProperty("action-name", "win.next-comic")
	nb.nextButton.AddAccelerator("activate", nb.accels, gdk.KEY_Right, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
	nb.box.Add(nb.nextButton)

	nb.newestButton, err = gtk.ButtonNew()
	if err != nil {
		return nil, err
	}
	nb.newestButton.SetTooltipText(l("Go to the newest comic"))
	nb.newestButton.SetProperty("action-name", "win.newest-comic")
	nb.newestButton.AddAccelerator("activate", nb.accels, gdk.KEY_End, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
	nb.box.Add(nb.newestButton)

	return nb, nil
}

func (nb *NavigationBar) Destroy() {
	nb.accels = nil

	nb.box = nil

	nb.firstButton = nil
	nb.previousButton = nil
	nb.randomButton = nil
	nb.nextButton = nil
	nb.newestButton = nil
}

func (nb *NavigationBar) IWidget() gtk.IWidget {
	return nb.box
}

func (nb *NavigationBar) SetFirstButtonImage(image gtk.IWidget) {
	nb.firstButton.SetImage(image)
}

func (nb *NavigationBar) SetPreviousButtonImage(image gtk.IWidget) {
	nb.previousButton.SetImage(image)
}

func (nb *NavigationBar) SetRandomButtonImage(image gtk.IWidget) {
	nb.randomButton.SetImage(image)
}

func (nb *NavigationBar) SetNextButtonImage(image gtk.IWidget) {
	nb.nextButton.SetImage(image)
}

func (nb *NavigationBar) SetNewestButtonImage(image gtk.IWidget) {
	nb.newestButton.SetImage(image)
}

func (nb *NavigationBar) SetLinkedButtons(linked bool) {
	sc, err := nb.box.GetStyleContext()
	if err != nil {
		log.Print(err)
		return
	}

	if linked {
		sc.AddClass(style.ClassLinked)
		nb.box.SetSpacing(0)
	} else {
		sc.RemoveClass(style.ClassLinked)
		nb.box.SetSpacing(4)
	}
}
