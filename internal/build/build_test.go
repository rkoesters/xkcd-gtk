package build

import (
	"testing"
)

func TestParse(t *testing.T) {
	const data = "version=0.0.0,debug=true"
	flags := parse(data)
	if flags == nil {
		t.Fatal("flags == nil")
	}
	if flags["version"] != "0.0.0" || flags["debug"] != "true" {
		t.Fatalf("parse failure: data='%v' flags='%v'", data, flags)
	}
}

func TestParseEmpty(t *testing.T) {
	const data = ""
	flags := parse(data)
	if flags == nil {
		t.Fatal("flags == nil")
	}
	if flags["version"] != "" || flags["debug"] != "" {
		t.Fatalf("parse failure: data='%v' flags='%v'", data, flags)
	}
}

func TestParseBlankKVPair(t *testing.T) {
	const data = "version=0.0.0,"
	flags := parse(data)
	if flags == nil {
		t.Fatal("flags == nil")
	}
	if flags["version"] != "0.0.0" || flags["debug"] != "" {
		t.Fatalf("parse failure: data='%v' flags='%v'", data, flags)
	}
	if flags[""] != "" {
		t.Fatal("parse failure: blank key=value pair has non-zero value")
	}
}

func TestParseBadFormat(t *testing.T) {
	// Should print warnings, but should not panic.
	const data = "asdf"
	flags := parse(data)
	if flags == nil {
		t.Fatal("flags == nil")
	}
	if flags["version"] != "" || flags["debug"] != "" {
		t.Fatalf("parse failure: data='%v' flags='%v'", data, flags)
	}
	if flags["asdf"] != "" {
		t.Fatal("parse failure: invalid key=value pair has non-zero value")
	}
}
