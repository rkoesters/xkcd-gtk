package widget

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"strconv"
	"strings"
	"time"
)

// PropertiesDialog holds a gtk dialog that shows the comic information for the
// parent window's comic.
type PropertiesDialog struct {
	parent *ApplicationWindow
	dialog *gtk.Dialog

	labels map[string]*gtk.Label

	accels *gtk.AccelGroup
}

var _ Window = &PropertiesDialog{}

const (
	propertiesKeyNumber     = "number"
	propertiesKeyTitle      = "title"
	propertiesKeyDate       = "date"
	propertiesKeyImage      = "image"
	propertiesKeyAltText    = "alt text"
	propertiesKeyNews       = "news"
	propertiesKeyLink       = "link"
	propertiesKeyTranscript = "transcript"
)

// NewPropertiesDialog creates and returns a PropertiesDialog for the given
// parent Window.
func NewPropertiesDialog(parent *ApplicationWindow) (*PropertiesDialog, error) {
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
	pd.accels.Connect(gdk.KEY_q, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE, parent.app.Quit)

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

	pd.addRowToGrid(grid, 0, propertiesKeyNumber, l("Number"))
	pd.addRowToGrid(grid, 1, propertiesKeyTitle, l("Title"))
	pd.addRowToGrid(grid, 2, propertiesKeyDate, l("Date"))
	pd.addRowToGrid(grid, 3, propertiesKeyImage, l("Image"))
	pd.addRowToGrid(grid, 4, propertiesKeyAltText, l("Alt Text"))
	pd.addRowToGrid(grid, 5, propertiesKeyNews, l("News"))
	pd.addRowToGrid(grid, 6, propertiesKeyLink, l("Link"))
	pd.addRowToGrid(grid, 7, propertiesKeyTranscript, l("Transcript"))
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
func (win *ApplicationWindow) ShowProperties() {
	var err error
	if win.properties == nil {
		win.properties, err = NewPropertiesDialog(win)
		if err != nil {
			log.Print("error creating properties dialog: ", err)
			return
		}
	}

	win.app.application.AddWindow(win.properties.IWindow())
	win.properties.dialog.Present()
}

func (pd *PropertiesDialog) addRowToGrid(grid *gtk.Grid, row int, key, label string) error {
	keyLabel, err := gtk.LabelNew(label)
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

	pd.labels[propertiesKeyNumber].SetText(strconv.Itoa(pd.parent.comic.Num))
	pd.labels[propertiesKeyTitle].SetText(pd.parent.comic.Title)
	pd.labels[propertiesKeyImage].SetText(pd.parent.comic.Img)
	pd.labels[propertiesKeyAltText].SetText(pd.parent.comic.Alt)
	pd.labels[propertiesKeyDate].SetText(formatDate(pd.parent.comic.Year, pd.parent.comic.Month, pd.parent.comic.Day))
	pd.labels[propertiesKeyNews].SetText(pd.parent.comic.News)
	pd.labels[propertiesKeyLink].SetText(pd.parent.comic.Link)
	pd.labels[propertiesKeyTranscript].SetText(pd.parent.comic.Transcript)
}

// Close is called when the dialog is closed. It tells the parent to save its
// window state.
func (pd *PropertiesDialog) Close() {
	pd.parent.properties = nil
	pd.parent.state.PropertiesWidth, pd.parent.state.PropertiesHeight = pd.dialog.GetSize()
	pd.parent.state.PropertiesPositionX, pd.parent.state.PropertiesPositionY = pd.dialog.GetPosition()
	pd.parent.state.SaveState(pd.parent.window, pd.dialog)
}

// Destroy removes our references to the dialog so the garbage collector can
// take care of it.
func (pd *PropertiesDialog) Destroy() {
	pd.parent = nil
	pd.dialog = nil

	pd.labels = nil

	pd.accels = nil
}

func (pd *PropertiesDialog) IWidget() gtk.IWidget { return pd.dialog }
func (pd *PropertiesDialog) IWindow() gtk.IWindow { return pd.dialog }

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
