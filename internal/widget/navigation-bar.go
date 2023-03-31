package widget

import (
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

type NavigationBar struct {
	*gtk.ButtonBox

	firstButton    *gtk.Button
	previousButton *gtk.Button
	randomButton   *gtk.Button
	nextButton     *gtk.Button
	newestButton   *gtk.Button

	// For UpdateButtonState.
	actions     map[string]*glib.SimpleAction
	comicNumber func() int
}

var _ Widget = &NavigationBar{}

func NewNavigationBar(accels *gtk.AccelGroup, actions map[string]*glib.SimpleAction, comicNumber func() int) (*NavigationBar, error) {
	super, err := gtk.ButtonBoxNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	nb := &NavigationBar{
		ButtonBox: super,

		actions:     actions,
		comicNumber: comicNumber,
	}

	nb.SetLayout(gtk.BUTTONBOX_EXPAND)

	addNavButton := func(label, action string, key uint) (*gtk.Button, error) {
		btn, err := gtk.ButtonNew()
		if err != nil {
			return nil, err
		}
		btn.SetTooltipText(label)
		btn.SetActionName(action)
		btn.AddAccelerator("activate", accels, key, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE)
		nb.Add(btn)
		return btn, nil
	}

	nb.firstButton, err = addNavButton(l("Go to the first comic"), "win.first-comic", gdk.KEY_Home)
	if err != nil {
		return nil, err
	}

	nb.previousButton, err = addNavButton(l("Go to the previous comic"), "win.previous-comic", gdk.KEY_Left)
	if err != nil {
		return nil, err
	}

	nb.randomButton, err = addNavButton(l("Go to a random comic"), "win.random-comic", gdk.KEY_r)
	if err != nil {
		return nil, err
	}

	nb.nextButton, err = addNavButton(l("Go to the next comic"), "win.next-comic", gdk.KEY_Right)
	if err != nil {
		return nil, err
	}

	nb.newestButton, err = addNavButton(l("Go to the newest comic"), "win.newest-comic", gdk.KEY_End)
	if err != nil {
		return nil, err
	}

	return nb, nil
}

func (nb *NavigationBar) Dispose() {
	if nb == nil {
		return
	}

	nb.ButtonBox = nil

	nb.firstButton = nil
	nb.previousButton = nil
	nb.randomButton = nil
	nb.nextButton = nil
	nb.newestButton = nil
}

func (nb *NavigationBar) UpdateButtonState(comicNumber int) {
	nb.actions["first-comic"].SetEnabled(comicNumber > 1)
	nb.actions["previous-comic"].SetEnabled(comicNumber > 1)

	newest, _ := cache.NewestComicInfoFromCache()
	nb.actions["next-comic"].SetEnabled(comicNumber < newest.Num)

	go func() {
		newest, _ := cache.CheckForNewestComicInfo(time.Hour)
		glib.IdleAddPriority(glib.PRIORITY_DEFAULT, func() {
			n := nb.comicNumber()
			nb.actions["next-comic"].SetEnabled(n < newest.Num)
		})
	}()
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

func (nb *NavigationBar) SetLinkedButtons(linked bool) error {
	sc, err := nb.GetStyleContext()
	if err != nil {
		return err
	}

	if linked {
		sc.AddClass(style.ClassLinked)
		nb.SetSpacing(0)
	} else {
		sc.RemoveClass(style.ClassLinked)
		nb.SetSpacing(style.PaddingUnlinkedButtonBox)
	}

	return nil
}
