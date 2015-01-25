package buffer

import (
	"encoding/gob"
	"io"
	"io/ioutil"
)

type discard struct{}

// NewDiscard creates a Buffer which writes to ioutil.Discard and read's return 0, io.EOF.
func NewDiscard() Buffer {
	return &discard{}
}

func (buf *discard) Len() int64 {
	return 0
}

func (buf *discard) Cap() int64 {
	return MAXINT64
}

func (buf *discard) Reset() {}

func (buf *discard) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (buf *discard) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, io.EOF
}

func (buf *discard) Write(p []byte) (int, error) {
	return ioutil.Discard.Write(p)
}

func (buf *discard) WriteAt(p []byte, off int64) (int, error) {
	return ioutil.Discard.Write(p)
}

func init() {
	gob.Register(&discard{})
}
