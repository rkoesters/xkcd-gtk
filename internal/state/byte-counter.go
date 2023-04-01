package state

import (
	"io"
)

type byteCounter struct {
	io.Writer
	io.Reader
	count int64
}

var _ io.ReadWriter = &byteCounter{}

func (bc *byteCounter) Write(p []byte) (n int, err error) {
	n, err = bc.Writer.Write(p)
	bc.count = bc.count + int64(n)
	return
}

func (bc *byteCounter) Read(p []byte) (n int, err error) {
	n, err = bc.Reader.Read(p)
	bc.count = bc.count + int64(n)
	return
}
