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

const (
	whatIfLink = "https://what-if.xkcd.com/"
	blogLink   = "https://blog.xkcd.com/"
	storeLink  = "https://store.xkcd.com/"
	aboutLink  = "https://xkcd.com/about/"
)

// OpenWhatIf opens whatifLink in the user's web browser.
func (app *Application) OpenWhatIf() {
	app.OpenURL(whatIfLink)
}

// OpenBlog opens blogLink in the user's web browser.
func (app *Application) OpenBlog() {
	app.OpenURL(blogLink)
}

// OpenStore opens storeLink in the user's web browser.
func (app *Application) OpenStore() {
	app.OpenURL(storeLink)
}

// OpenAboutXKCD opens aboutLink in the user's web browser.
func (app *Application) OpenAboutXKCD() {
	app.OpenURL(aboutLink)
}
