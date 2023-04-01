package state

import (
	"io"
)

type byteCounter struct {
	io.Writer
	count int64
}

var _ io.Writer = &byteCounter{}

func (bc *byteCounter) Write(p []byte) (n int, err error) {
	n, err = bc.Writer.Write(p)
	bc.count = bc.count + int64(n)
	return
}
