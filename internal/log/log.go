// Package log provides functions for logging fatal errors, warnings, and debug
// info.
package log

import (
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
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
}
