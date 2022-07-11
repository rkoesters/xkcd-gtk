// Package build provides access to compile-time build information.
package build

import (
	"github.com/rkoesters/xkcd-gtk/internal/log"
	"strings"
)

// Comma separated list of key=value pairs.
var data = ""

// Options is a [key]=value pair of compile time flags. Init must be called
// before using Options.
var Options map[string]string

// Init initializes the build package by parsing the data provided to it at
// compile time. Init must be called before using Options or calling any other
// function provided by this package.
func Init() {
	log.Debug("build data: ", data)

	Options = parse(data)
}

func parse(data string) map[string]string {
	flags := make(map[string]string)

	if data == "" {
		return flags
	}

	for _, s := range strings.Split(data, ",") {
		pair := strings.SplitN(s, "=", 2)
		if len(pair) != 2 {
			log.Print("error parsing build flag: ", s)
			continue
		}
		flags[pair[0]] = pair[1]
	}

	return flags
}

// Version returns the version string of this binary.
func Version() string {
	v, ok := Options["version"]
	if ok {
		return v
	} else {
		return "undefined"
	}
}
