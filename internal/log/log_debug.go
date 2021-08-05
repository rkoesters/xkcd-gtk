// +build xkcd_gtk_debug

package log

import (
	"log"
)

var (
	Debug   = log.Print
	Debugf  = log.Printf
	Debugln = log.Println
)
