package state_test

import (
	"bytes"
	"testing"

	"github.com/rkoesters/xkcd-gtk/internal/state"
)

func TestWriteTo(t *testing.T) {
	tests := []struct {
		name string
		n    int64 // len(json(ws)) + len('\n')
		ws   state.Window
		json string
	}{{
		name: "empty state",
		n:    215,
		ws:   state.Window{},
		json: `{"ComicNumber":0,"Maximized":false,"Height":0,"Width":0,"PositionX":0,"PositionY":0,"ImageScale":0,"PropertiesVisible":false,"PropertiesHeight":0,"PropertiesWidth":0,"PropertiesPositionX":0,"PropertiesPositionY":0}
`,
	}, {
		name: "maximized",
		n:    224,
		ws: state.Window{
			ComicNumber: 123,
			Maximized:   true,
			Height:      1080,
			Width:       1920,
			ImageScale:  1.5,
		},
		json: `{"ComicNumber":123,"Maximized":true,"Height":1080,"Width":1920,"PositionX":0,"PositionY":0,"ImageScale":1.5,"PropertiesVisible":false,"PropertiesHeight":0,"PropertiesWidth":0,"PropertiesPositionX":0,"PropertiesPositionY":0}
`,
	}, {
		name: "with properties",
		n:    224,
		ws: state.Window{
			ComicNumber:       123,
			Height:            500,
			Width:             700,
			ImageScale:        1,
			PropertiesVisible: true,
			PropertiesHeight:  400,
			PropertiesWidth:   300,
		},
		json: `{"ComicNumber":123,"Maximized":false,"Height":500,"Width":700,"PositionX":0,"PositionY":0,"ImageScale":1,"PropertiesVisible":true,"PropertiesHeight":400,"PropertiesWidth":300,"PropertiesPositionX":0,"PropertiesPositionY":0}
`,
	}}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			n, err := tc.ws.WriteTo(&buf)
			if err != nil {
				t.Error("WriteTo returned error: ", err)
			}
			if n != tc.n {
				t.Errorf("wrote %v bytes, want %v", n, tc.n)
			}
			s := buf.String()
			if s != tc.json {
				t.Errorf("wrong json, got=%q, want=%q", s, tc.json)
			}
		})
	}
}
