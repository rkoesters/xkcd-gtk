package log_test

import (
	"io/ioutil"
	"log"
	"testing"
)

// This log package makes heavy use of the internal/build.Debug() function,
// which was originally implemented as a map lookup and string comparison. These
// benchmarks demonstrate the performance difference between using the
// map-lookup and string-comparison vs caching the bool result of the
// comparison.

var logger *log.Logger

func init() {
	logger = log.New(ioutil.Discard, "", log.LstdFlags|log.Lshortfile)
}

func BenchmarkLogDebugMapOff(b *testing.B) {
	options := make(map[string]string)

	logDebug := func(v ...interface{}) {
		if options["debug"] == "on" {
			logger.Print(v...)
		}
	}

	for i := 0; i < b.N; i++ {
		logDebug("test log message")
	}
}

func BenchmarkLogDebugMapOn(b *testing.B) {
	options := make(map[string]string)
	options["debug"] = "on"

	logDebug := func(v ...interface{}) {
		if options["debug"] == "on" {
			logger.Print(v...)
		}
	}

	for i := 0; i < b.N; i++ {
		logDebug("test log message")
	}
}

func BenchmarkLogDebugBoolOff(b *testing.B) {
	options := make(map[string]string)

	debug := options["debug"] == "on"

	logDebug := func(v ...interface{}) {
		if debug {
			logger.Print(v...)
		}
	}

	for i := 0; i < b.N; i++ {
		logDebug("test log message")
	}
}

func BenchmarkLogDebugBoolOn(b *testing.B) {
	options := make(map[string]string)
	options["debug"] = "on"

	debug := options["debug"] == "on"

	logDebug := func(v ...interface{}) {
		if debug {
			logger.Print(v...)
		}
	}

	for i := 0; i < b.N; i++ {
		logDebug("test log message")
	}
}
