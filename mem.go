package buffer

import (
	"bytes"
	"io"
)

type MemBuffer struct {
	limit int64
	*bytes.Buffer
}

func New(n int64) *MemBuffer {
	return &MemBuffer{
		limit:  n,
		Buffer: bytes.NewBuffer(nil),
	}
}

func (buf *MemBuffer) Cap() int64 {
	return buf.limit
}

func (buf *MemBuffer) Len() int64 {
	return int64(buf.Buffer.Len())
}

func (buf *MemBuffer) Write(p []byte) (n int, err error) {
	return LimitWriter(buf.Buffer, Gap(buf)).Write(p)
}

func (buf *MemBuffer) ReadAt(b []byte, off int64) (n int, err error) {
	return bytes.NewReader(buf.Bytes()).ReadAt(b, off)
}

func (buf *MemBuffer) Read(p []byte) (n int, err error) {
	return io.LimitReader(buf.Buffer, buf.Len()).Read(p)
}

func (buf *MemBuffer) FastForward(n int) int {
	data := buf.Bytes()
	if n > len(data) {
		n = len(data)
	}
	b := bytes.NewBuffer(data[n:])
	buf.Buffer = b
	return n
}
