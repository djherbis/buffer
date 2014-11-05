package buffer

import (
	"io"
	"io/ioutil"
)

type Discard struct {
	io.Writer
}

func NewDiscard() *Discard {
	return &Discard{Writer: ioutil.Discard}
}

func (buf *Discard) Len() int64 {
	return 0
}

func (buf *Discard) Cap() int64 {
	return 0
}

func (buf *Discard) Reset() {}

func (buf *Discard) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (buf *Discard) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, io.EOF
}
