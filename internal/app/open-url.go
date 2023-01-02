package app

import (
	"github.com/rkoesters/xdg"
	"github.com/rkoesters/xkcd-gtk/internal/log"
)

// OpenURL opens the given URL in the user's web browser. The function logs
// errors in addition to returning them, so errors can be ignored if that is all
// that would be done with the returned error.
func (app *Application) OpenURL(url string) error {
	err := xdg.Open(url)
	if err != nil {
		log.Print("error opening ", url, " in web browser: ", err)
	}
	return err
}