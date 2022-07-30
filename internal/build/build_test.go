package build

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		desc string
		data string
		want map[string]string
	}{{
		desc: "empty build data",
		want: map[string]string{},
	}, {
		desc: "well formed build data",
		data: "version=0.0.0,debug=true",
		want: map[string]string{
			"version": "0.0.0",
			"debug":   "true",
		},
	}, {
		desc: "blank key-value pair",
		data: "version=0.0.0,",
		want: map[string]string{
			"version": "0.0.0",
			"":        "",
		},
	}, {
		desc: "bad format",
		data: "asdf",
		want: map[string]string{
			"asdf": "",
		},
	}}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			flags := parse(tc.data)
			if flags == nil {
				t.Fatal("flags == nil")
			}
			for k, got := range flags {
				want, ok := tc.want[k]
				if !ok {
					t.Errorf("unexpected key %q", k)
				}
				if got != want {
					t.Errorf("got %q, want %q", got, want)
				}
				delete(tc.want, k)
			}
			if len(tc.want) > 0 {
				t.Error("missing key-value pairs: ", tc.want)
			}
		})
	}
}
