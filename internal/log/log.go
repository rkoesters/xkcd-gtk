// Package log provides functions for logging fatal errors, warnings, and debug
// info.
package log

import (
	"fmt"
	"github.com/rkoesters/xkcd-gtk/internal/build"
	"log"
)

var (
	Fatal   = log.Fatal
	Fatalf  = log.Fatalf
	Fatalln = log.Fatalln

	Panic   = log.Panic
	Panicf  = log.Panicf
	Panicln = log.Panicln

	Print   = log.Print
	Printf  = log.Printf
	Println = log.Println
)

// Init performs initialization required for the log package.
func Init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Debug calls log.Print if Debug returns true.
func Debug(v ...interface{}) {
	if build.Debug() {
		log.Output(2, fmt.Sprint(v...))
	}
}

// Debugf calls log.Printf if Debug returns true.
func Debugf(format string, v ...interface{}) {
	if build.Debug() {
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

// Debugln calls log.Println if Debug returns true.
func Debugln(v ...interface{}) {
	if build.Debug() {
		log.Output(2, fmt.Sprintln(v...))
	}
}
