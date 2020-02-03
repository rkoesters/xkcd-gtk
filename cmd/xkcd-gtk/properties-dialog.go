package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"strconv"
	"strings"
	"time"
)

// PropertiesDialog holds a gtk dialog that shows the comic information for the
// parent window's comic.
type PropertiesDialog struct {
	parent *Window
	dialog *gtk.Dialog

	labels map[string]*gtk.Label

	accels *gtk.AccelGroup
}

// NewPropertiesDialog creates and returns a PropertiesDialog for the given
// parent Window.
func NewPropertiesDialog(parent *Window) (*PropertiesDialog, error) {
	var err error

	pd := new(PropertiesDialog)
	pd.labels = make(map[string]*gtk.Label)
	pd.parent = parent

	pd.dialog, err = gtk.DialogNew()
	if err != nil {
		return nil, err
	}
	pd.dialog.SetTransientFor(parent.window)
	pd.dialog.SetTitle(l("Properties"))
	pd.dialog.SetSizeRequest(400, 400)
	pd.dialog.SetDestroyWithParent(true)
	pd.dialog.Resize(parent.state.PropertiesWidth, parent.state.PropertiesHeight)
	if parent.state.PropertiesPositionX != 0 && parent.state.PropertiesPositionY != 0 {
		pd.dialog.Move(parent.state.PropertiesPositionX, parent.state.PropertiesPositionY)
	}

	// Make Control-q to quit the app work in this dialog.
	pd.accels, err = gtk.AccelGroupNew()
	if err != nil {
		return nil, err
	}
	pd.dialog.AddAccelGroup(pd.accels)
	pd.accels.Connect(gdk.KEY_q, gdk.GDK_CONTROL_MASK, gtk.ACCEL_VISIBLE, parent.app.Quit)

	pd.dialog.Connect("delete-event", pd.Close)
	pd.dialog.Connect("destroy", pd.Destroy)

	scwin, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	scwin.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	scwin.SetMarginTop(0)
	scwin.SetMarginBottom(0)
	scwin.SetMarginStart(0)
	scwin.SetMarginEnd(0)
	scwin.SetVExpand(true)

	grid, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}
	grid.SetColumnSpacing(12)
	grid.SetRowSpacing(12)
	grid.SetMarginTop(12)
	grid.SetMarginBottom(12)
	grid.SetMarginStart(12)
	grid.SetMarginEnd(12)

	pd.addRowToGrid(grid, 0, l("Number"))
	pd.addRowToGrid(grid, 1, l("Title"))
	pd.addRowToGrid(grid, 2, l("Date"))
	pd.addRowToGrid(grid, 3, l("Image"))
	pd.addRowToGrid(grid, 4, l("Alt Text"))
	pd.addRowToGrid(grid, 5, l("News"))
	pd.addRowToGrid(grid, 6, l("Link"))
	pd.addRowToGrid(grid, 7, l("Transcript"))
	pd.Update()

	scwin.Add(grid)

	box, err := pd.dialog.GetContentArea()
	if err != nil {
		return nil, err
	}
	// A gtk.Dialog content area has some children by default, we want to
	// remove those children so the only child is scwin.
	box.GetChildren().Foreach(func(child interface{}) {
		box.Remove(child.(gtk.IWidget))
	})
	box.Add(scwin)
	box.ShowAll()

	return pd, nil
}

// ShowProperties presents the properties dialog to the user. If the dialog
// doesn't exist yet, we create it.
func (win *Window) ShowProperties() {
	var err error
	if win.properties == nil {
		win.properties, err = NewPropertiesDialog(win)
		if err != nil {
			log.Print("error creating properties dialog: ", err)
			return
		}
	}

	win.app.application.AddWindow(&win.properties.dialog.Window)
	win.properties.dialog.Present()
}

func (pd *PropertiesDialog) addRowToGrid(grid *gtk.Grid, row int, key string) error {
	keyLabel, err := gtk.LabelNew(key)
	if err != nil {
		return err
	}
	keyLabel.SetHAlign(gtk.ALIGN_END)
	keyLabel.SetVAlign(gtk.ALIGN_START)
	valLabel, err := gtk.LabelNew("")
	if err != nil {
		return err
	}
	valLabel.SetXAlign(0)
	valLabel.SetHAlign(gtk.ALIGN_START)
	valLabel.SetVAlign(gtk.ALIGN_START)
	valLabel.SetLineWrap(true)
	valLabel.SetSelectable(true)
	valLabel.SetCanFocus(false)

	grid.Attach(keyLabel, 0, row, 1, 1)
	grid.Attach(valLabel, 1, row, 1, 1)

	pd.labels[key] = valLabel

	return nil
}

// Update changes the dialog's contents to match the parent Window's comic.
func (pd *PropertiesDialog) Update() {
	pd.parent.comicMutex.RLock()
	defer pd.parent.comicMutex.RUnlock()

	pd.labels[l("Number")].SetText(strconv.Itoa(pd.parent.comic.Num))
	pd.labels[l("Title")].SetText(pd.parent.comic.Title)
	pd.labels[l("Image")].SetText(pd.parent.comic.Img)
	pd.labels[l("Alt Text")].SetText(pd.parent.comic.Alt)
	pd.labels[l("Date")].SetText(formatDate(pd.parent.comic.Year, pd.parent.comic.Month, pd.parent.comic.Day))
	pd.labels[l("News")].SetText(pd.parent.comic.News)
	pd.labels[l("Link")].SetText(pd.parent.comic.Link)
	pd.labels[l("Transcript")].SetText(pd.parent.comic.Transcript)
}

// Close is called when the dialog is closed. It tells the parent to save its
// window state.
func (pd *PropertiesDialog) Close() {
	pd.parent.properties = nil
	pd.parent.state.PropertiesWidth, pd.parent.state.PropertiesHeight = pd.dialog.GetSize()
	pd.parent.state.PropertiesPositionX, pd.parent.state.PropertiesPositionY = pd.dialog.GetPosition()
	pd.parent.SaveState()
}

// Destroy removes our references to the dialog so the garbage collector can
// take care of it.
func (pd *PropertiesDialog) Destroy() {
	pd.labels = nil
	pd.dialog = nil
	pd.parent = nil
}

// formatDate takes a year, month, and date as strings and turns them into a
// pretty date.
func formatDate(year, month, day string) string {
	date := strings.Join([]string{year, month, day}, "-")
	t, err := time.Parse("2006-1-2", date)
	if err != nil {
		return ""
	}
	return t.Format("Jan _2, 2006")
}
