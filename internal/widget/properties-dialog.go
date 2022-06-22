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

	comicNumber     *gtk.Label
	comicTitle      *gtk.Label
	comicDate       *gtk.Label
	comicImage      *gtk.Label
	comicAltText    *gtk.Label
	comicNews       *gtk.Label
	comicLink       *gtk.Label
	comicTranscript *gtk.Label
}

var _ Window = &PropertiesDialog{}

// NewPropertiesDialog creates and returns a PropertiesDialog for the given
// parent Window.
func NewPropertiesDialog(parent *ApplicationWindow) (*PropertiesDialog, error) {
	var err error

	pd := new(PropertiesDialog)
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
	accels, err := gtk.AccelGroupNew()
	if err != nil {
		return nil, err
	}
	pd.dialog.AddAccelGroup(accels)
	accels.Connect(gdk.KEY_q, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE, parent.app.Quit)

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

	row := 0

	addRowToGrid := func(label string) (*gtk.Label, error) {
		defer func() { row++ }()

		keyLabel, err := gtk.LabelNew(label)
		if err != nil {
			return nil, err
		}
		keyLabel.SetHAlign(gtk.ALIGN_END)
		keyLabel.SetVAlign(gtk.ALIGN_START)
		valLabel, err := gtk.LabelNew("")
		if err != nil {
			return nil, err
		}
		valLabel.SetXAlign(0)
		valLabel.SetHAlign(gtk.ALIGN_START)
		valLabel.SetVAlign(gtk.ALIGN_START)
		valLabel.SetLineWrap(true)
		valLabel.SetSelectable(true)
		valLabel.SetCanFocus(false)

		grid.Attach(keyLabel, 0, row, 1, 1)
		grid.Attach(valLabel, 1, row, 1, 1)

		return valLabel, nil
	}

	pd.comicNumber, err = addRowToGrid(l("Number"))
	if err != nil {
		return nil, err
	}
	pd.comicTitle, err = addRowToGrid(l("Title"))
	if err != nil {
		return nil, err
	}
	pd.comicDate, err = addRowToGrid(l("Date"))
	if err != nil {
		return nil, err
	}
	pd.comicImage, err = addRowToGrid(l("Image"))
	if err != nil {
		return nil, err
	}
	pd.comicAltText, err = addRowToGrid(l("Alt text"))
	if err != nil {
		return nil, err
	}
	pd.comicNews, err = addRowToGrid(l("News"))
	if err != nil {
		return nil, err
	}
	pd.comicLink, err = addRowToGrid(l("Link"))
	if err != nil {
		return nil, err
	}
	pd.comicTranscript, err = addRowToGrid(l("Transcript"))
	if err != nil {
		return nil, err
	}
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

// Update changes the dialog's contents to match the parent Window's comic.
func (pd *PropertiesDialog) Update() {
	pd.parent.comicMutex.RLock()
	defer pd.parent.comicMutex.RUnlock()

	pd.comicNumber.SetText(strconv.Itoa(pd.parent.comic.Num))
	pd.comicTitle.SetText(pd.parent.comic.Title)
	pd.comicImage.SetText(pd.parent.comic.Img)
	pd.comicAltText.SetText(pd.parent.comic.Alt)
	pd.comicDate.SetText(formatDate(pd.parent.comic.Year, pd.parent.comic.Month, pd.parent.comic.Day))
	pd.comicNews.SetText(pd.parent.comic.News)
	pd.comicLink.SetText(pd.parent.comic.Link)
	pd.comicTranscript.SetText(pd.parent.comic.Transcript)
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

	pd.comicNumber = nil
	pd.comicTitle = nil
	pd.comicImage = nil
	pd.comicAltText = nil
	pd.comicDate = nil
	pd.comicNews = nil
	pd.comicLink = nil
	pd.comicTranscript = nil
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
