//go:build !xkcd_gtk_debug

package log

// Debug is equivalent to Print if this is a debug build, otherwise it is a
// no-op.
func Debug(v ...any) {}

// Debugf is equivalent to Printf if this is a debug build, otherwise it is a
// no-op.
func Debugf(format string, v ...any) {}

// Debugln is equivalent to Println if this is a debug build, otherwise it is a
// no-op.
func Debugln(v ...any) {}
