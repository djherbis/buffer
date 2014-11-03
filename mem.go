package buffer

import (
	"bytes"
	"io"
)

type MemBuffer struct {
	limit int64
	bytes.Buffer
}

func New(n int64) *MemBuffer {
	return &MemBuffer{
		limit: n,
	}
}

func (buf *MemBuffer) Cap() int64 {
	return buf.limit
}

func (buf *MemBuffer) Len() int64 {
	return int64(buf.Buffer.Len())
}

func (buf *MemBuffer) Write(p []byte) (n int, err error) {
	n, err = buf.Buffer.Write(ShrinkToFit(buf, p))
	return n, err
}

func (buf *MemBuffer) ReadAt(b []byte, off int64) (n int, err error) {
	return bytes.NewReader(buf.Bytes()).ReadAt(b, off)
}

func (buf *MemBuffer) Read(p []byte) (n int, err error) {
	n, err = buf.Buffer.Read(ShrinkToRead(buf, p))
	if buf.Len() == 0 {
		err = io.EOF
	}
	return n, err
}
