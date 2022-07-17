package widget

import (
	"strconv"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"github.com/rkoesters/xkcd-gtk/internal/style"
)

// PropertiesDialog holds a gtk dialog that shows the comic information for the
// parent window's comic.
type PropertiesDialog struct {
	*gtk.Dialog

	parent *ApplicationWindow

	comicNumber     *gtk.Label
	comicTitle      *gtk.Label
	comicDate       *gtk.Label
	comicImage      *gtk.Label
	comicAltText    *gtk.Label
	comicNews       *gtk.Label
	comicLink       *gtk.Label
	comicTranscript *gtk.Label
}

var _ Widget = &PropertiesDialog{}

// NewPropertiesDialog creates and returns a PropertiesDialog for the given
// parent Window.
func NewPropertiesDialog(parent *ApplicationWindow) (*PropertiesDialog, error) {
	super, err := gtk.DialogNew()
	if err != nil {
		return nil, err
	}
	pd := &PropertiesDialog{
		Dialog: super,

		parent: parent,
	}

	pd.SetTransientFor(parent.ApplicationWindow)
	pd.SetTitle(l("Properties"))
	pd.SetSizeRequest(400, 400)
	pd.SetDestroyWithParent(true)
	pd.Resize(parent.state.PropertiesWidth, parent.state.PropertiesHeight)
	if parent.state.PropertiesPositionX != 0 && parent.state.PropertiesPositionY != 0 {
		pd.Move(parent.state.PropertiesPositionX, parent.state.PropertiesPositionY)
	}

	// Initialize our window accelerators.
	accels, err := gtk.AccelGroupNew()
	if err != nil {
		return nil, err
	}
	pd.AddAccelGroup(accels)
	accels.Connect(gdk.KEY_w, gdk.CONTROL_MASK, gtk.ACCEL_VISIBLE, pd.Close)

	pd.Connect("delete-event", pd.DeleteEvent)
	pd.Connect("destroy", pd.Dispose)

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
	grid.SetColumnSpacing(style.PaddingPropertiesDialog)
	grid.SetRowSpacing(style.PaddingPropertiesDialog)
	grid.SetMarginTop(style.PaddingPropertiesDialog)
	grid.SetMarginBottom(style.PaddingPropertiesDialog)
	grid.SetMarginStart(style.PaddingPropertiesDialog)
	grid.SetMarginEnd(style.PaddingPropertiesDialog)

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

	box, err := pd.GetContentArea()
	if err != nil {
		return nil, err
	}
	// A gtk.Dialog content area has some children by default, we want to
	// remove those children so the only child is scwin.
	box.GetChildren().Foreach(func(child interface{}) {
		w, ok := child.(*gtk.Widget)
		if !ok {
			log.Print("error converting child to gtk.Widget")
			return
		}
		box.Remove(w)
	})
	box.Add(scwin)
	box.ShowAll()

	return pd, nil
}

// ShowProperties presents the properties dialog to the user. If the dialog
// doesn't exist yet, we create it.
func (win *ApplicationWindow) ShowProperties() {
	if win.properties == nil {
		pd, err := NewPropertiesDialog(win)
		if err != nil {
			log.Print("error creating properties dialog: ", err)
			return
		}
		win.properties = pd
	}
	win.app.AddWindow(win.properties)
	win.properties.Dialog.Present()
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

// DeleteEvent is called when the dialog is closed. It tells the parent to save
// its window state.
func (pd *PropertiesDialog) DeleteEvent() {
	pd.parent.properties = nil
	pd.parent.state.PropertiesWidth, pd.parent.state.PropertiesHeight = pd.GetSize()
	pd.parent.state.PropertiesPositionX, pd.parent.state.PropertiesPositionY = pd.GetPosition()
	pd.parent.state.SaveState(pd.parent, pd)
}

// Dispose removes our references to the dialog so the garbage collector can
// take care of it.
func (pd *PropertiesDialog) Dispose() {
	if pd == nil {
		return
	}

	pd.Dialog = nil

	pd.parent = nil

	pd.comicNumber = nil
	pd.comicTitle = nil
	pd.comicImage = nil
	pd.comicAltText = nil
	pd.comicDate = nil
	pd.comicNews = nil
	pd.comicLink = nil
	pd.comicTranscript = nil
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
