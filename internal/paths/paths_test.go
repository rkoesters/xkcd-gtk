package paths_test

import (
	"github.com/rkoesters/xkcd-gtk/internal/paths"
	"path/filepath"
	"testing"
)

func TestPaths(t *testing.T) {
	// More thorough tests can be found in builder_test.go. These tests are
	// to verify that the package level versions of each method generally
	// works by checking that they return valid, absolute paths (or the
	// current directory in the case of the fallback for LocaleDir).

	paths.Init("com.example.test")

	if !filepath.IsAbs(paths.CacheDir()) {
		t.Fail()
	}
	if !filepath.IsAbs(paths.ConfigDir()) {
		t.Fail()
	}
	if !filepath.IsAbs(paths.DataDir()) {
		t.Fail()
	}
	if dir := paths.LocaleDir(); !filepath.IsAbs(dir) && dir != "." {
		t.Fail()
	}
}
