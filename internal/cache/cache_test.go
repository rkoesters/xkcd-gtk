package cache

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
)

func TestIntToBytes(t *testing.T) {
	ints := []int{-1, 0, 1, math.MaxInt32, math.MaxUint32, math.MinInt32}

	for _, n := range ints {
		b := intToBytes(n)
		t.Logf("%#v", b)

		buf := bytes.NewBuffer(b)
		nout, err := binary.ReadVarint(buf)
		if err != nil {
			t.Fatal(err)
		}

		if n != int(nout) {
			t.Fail()
		}
	}
}
