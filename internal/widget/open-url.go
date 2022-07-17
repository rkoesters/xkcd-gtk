package widget

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xdg"
	"github.com/rkoesters/xkcd-gtk/internal/log"
)

// openURL opens the given URL in the user's web browser. The function logs
// errors in addition to returning them, so errors can be ignored if that is all
// that would be done with the returned error.
func openURL(url string) error {
	err := xdg.Open(url)
	if err != nil {
		log.Print("error opening ", url, " in web browser: ", err)
	}
	return err
}

// urlLabel returns the given label string with an added suffix hinting to the
// user that the button will open a browser window.
func urlLabel(s string) string {
	switch gtk.GetLocaleDirection() {
	case gtk.TEXT_DIR_RTL:
		return "↖ " + s
	default:
		return s + " ➚"
	}
}
