package buffer

import (
	"io"

	"github.com/djherbis/buffer/wrapio"
)

type Ring struct {
	BufferAt
	L int64
	*wrapio.WrapReader
	*wrapio.WrapWriter
}

func NewRing(buffer BufferAt) Buffer {
	return &Ring{
		BufferAt:   buffer,
		WrapReader: wrapio.NewWrapReader(buffer, 0, buffer.Cap()),
		WrapWriter: wrapio.NewWrapWriter(buffer, 0, buffer.Cap()),
	}
}

func (buf *Ring) Len() int64 {
	return buf.L
}

func (buf *Ring) Cap() int64 {
	return MAXINT64
}

func (buf *Ring) Read(p []byte) (n int, err error) {
	if buf.Len() == buf.BufferAt.Cap() {
		buf.WrapReader.Seek(buf.WrapWriter.Offset(), 0)
	}
	n, err = io.LimitReader(buf.WrapReader, buf.Len()).Read(p)
	buf.L -= int64(n)
	return n, err
}

func (buf *Ring) Write(p []byte) (n int, err error) {
	n, err = buf.WrapWriter.Write(p)
	buf.L += int64(n)
	if buf.L > buf.BufferAt.Cap() {
		buf.L = buf.BufferAt.Cap()
	}
	return n, err
}

func (buf *Ring) Reset() {
	buf.BufferAt.Reset()
	buf.L = 0
	buf.WrapReader = wrapio.NewWrapReader(buf.BufferAt, 0, buf.BufferAt.Cap())
	buf.WrapWriter = wrapio.NewWrapWriter(buf.BufferAt, 0, buf.BufferAt.Cap())
}
