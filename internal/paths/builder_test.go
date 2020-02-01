package paths

import (
	"path/filepath"
	"strings"
	"testing"
)

const testAppID = "com.example.test"

func TestCacheDir(t *testing.T) {
	paths := builder{testAppID}

	dir := paths.CacheDir()

	if !filepath.IsAbs(dir) {
		t.Fail()
	}
	if !strings.HasSuffix(dir, testAppID) {
		t.Fail()
	}
}

func TestConfigDir(t *testing.T) {
	paths := builder{testAppID}

	dir := paths.ConfigDir()

	if !filepath.IsAbs(dir) {
		t.Fail()
	}
	if !strings.HasSuffix(dir, testAppID) {
		t.Fail()
	}
}

func TestDataDir(t *testing.T) {
	paths := builder{testAppID}

	dir := paths.DataDir()

	if !filepath.IsAbs(dir) {
		t.Fail()
	}
	if !strings.HasSuffix(dir, testAppID) {
		t.Fail()
	}
}

func TestLocaleDir(t *testing.T) {
	paths := builder{testAppID}

	dir := paths.LocaleDir()

	if !filepath.IsAbs(dir) && dir != "." {
		t.Fail()
	}
}
