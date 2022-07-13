package build_test

import (
	"testing"

	"github.com/gotk3/gotk3/glib"
	"github.com/rkoesters/xkcd-gtk/internal/build"
)

func TestAppIDIsValid(t *testing.T) {
	build.Init()
	if !glib.ApplicationIDIsValid(build.AppID()) {
		t.Error("invalid application ID:", build.AppID())
	}
}
