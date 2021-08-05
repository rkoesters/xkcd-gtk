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

func BenchmarkLogDebugOffMap(b *testing.B) {
	options := make(map[string]string)

	logDebug := func(v ...interface{}) {
		if options["debug"] == "on" {
			logger.Print(v...)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logDebug("test log message")
	}
}

func BenchmarkLogDebugOffBool(b *testing.B) {
	options := make(map[string]string)

	debug := options["debug"] == "on"

	logDebug := func(v ...interface{}) {
		if debug {
			logger.Print(v...)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logDebug("test log message")
	}
}

func BenchmarkLogDebugOffTag(b *testing.B) {
	options := make(map[string]string)
	_ = options["debug"] == "on"

	logDebug := func(v ...interface{}) {}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logDebug("test log message")
	}
}

func BenchmarkLogDebugOnMap(b *testing.B) {
	options := make(map[string]string)
	options["debug"] = "on"

	logDebug := func(v ...interface{}) {
		if options["debug"] == "on" {
			logger.Print(v...)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logDebug("test log message")
	}
}

func BenchmarkLogDebugOnBool(b *testing.B) {
	options := make(map[string]string)
	options["debug"] = "on"

	debug := options["debug"] == "on"

	logDebug := func(v ...interface{}) {
		if debug {
			logger.Print(v...)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logDebug("test log message")
	}
}

func BenchmarkLogDebugOnTag(b *testing.B) {
	options := make(map[string]string)
	options["debug"] = "on"

	logDebug := logger.Print

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logDebug("test log message")
	}
}

func BenchmarkLogDebugOnTagWithoutIndirection(b *testing.B) {
	options := make(map[string]string)
	options["debug"] = "on"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Print("test log message")
	}
}
