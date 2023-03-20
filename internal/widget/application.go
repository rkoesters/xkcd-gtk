package widget

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rkoesters/xkcd-gtk/internal/bookmarks"
	"github.com/rkoesters/xkcd-gtk/internal/cache"
)

// AppName is the user-visible name of this application.
func AppName() string { return l("Comic Sticks") }

// Application is interface needed by
// github.com/rkoesters/xkcd-gtk/internal/widget and implemented by
// github.com/rkoesters/xkcd-gtk/internal/app.
type Application interface {
	AddWindow(gtk.IWindow)
	BookmarksList() *bookmarks.List
	CacheWindowVR() cache.ViewRefresher
	CacheWindowVRW() cache.ViewRefreshWither
	ConnectDarkModeChanged(f interface{}) glib.SignalHandle
	DarkMode() bool
	GtkApplication() *gtk.Application
	GtkTheme() (string, error)
	OpenURL(string) error
	PrefersAppMenu() bool
	RemoveWindow(gtk.IWindow)
	SetDarkMode(bool)
}
