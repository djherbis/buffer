package buffer

import (
	"encoding/gob"
	"io"
	"io/ioutil"
)

type Discard struct{}

func NewDiscard() *Discard {
	return &Discard{}
}

func (buf *Discard) Len() int64 {
	return 0
}

func (buf *Discard) Cap() int64 {
	return MAXINT64
}

func (buf *Discard) Reset() {}

func (buf *Discard) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (buf *Discard) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, io.EOF
}

func (buf *Discard) Write(p []byte) (int, error) {
	return ioutil.Discard.Write(p)
}

func (buf *Discard) WriteAt(p []byte, off int64) (int, error) {
	return ioutil.Discard.Write(p)
}

func (buf *Discard) FFwd(n int64) int64 {
	return 0
}

func init() {
	gob.Register(&Discard{})
}
