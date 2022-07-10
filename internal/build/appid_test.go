package build_test

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/rkoesters/xkcd-gtk/internal/build"
	"testing"
)

func TestAppIDIsValid(t *testing.T) {
	if !glib.ApplicationIDIsValid(build.AppID) {
		t.Fail()
	}
}
