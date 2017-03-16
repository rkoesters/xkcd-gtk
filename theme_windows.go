// +build windows

package main

import (
	"os"
)

func init() {
	os.Setenv("GTK_THEME", "win32")
}
