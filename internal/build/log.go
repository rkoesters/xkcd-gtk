package build

import (
	"log"
)

// DebugPrint calls log.Print if Debug returns true.
func DebugPrint(v ...interface{}) {
	if Debug() {
		log.Print(v...)
	}
}

// DebugPrintf calls log.Printf if Debug returns true.
func DebugPrintf(format string, v ...interface{}) {
	if Debug() {
		log.Printf(format, v...)
	}
}

// DebugPrintln calls log.Println if Debug returns true.
func DebugPrintln(v ...interface{}) {
	if Debug() {
		log.Println(v...)
	}
}
