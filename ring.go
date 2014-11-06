package buffer

import "io"

type Ring struct {
	Buffer
	L int64
	*WrapReader
	*WrapWriter
}

func NewRing(buffer Buffer) *Ring {
	return &Ring{
		Buffer:     buffer,
		WrapReader: NewWrapReader(buffer, 0, buffer.Cap()),
		WrapWriter: NewWrapWriter(buffer, 0, buffer.Cap()),
	}
}

func (buf *Ring) Len() int64 {
	return buf.L
}

func (buf *Ring) Read(p []byte) (n int, err error) {
	if Full(buf) {
		buf.WrapReader.Seek(buf.WrapWriter.Offset(), 0)
	}
	n, err = io.LimitReader(buf.WrapReader, buf.Len()).Read(p)
	buf.L -= int64(n)
	return n, err
}

func (buf *Ring) ReadAt(p []byte, off int64) (n int, err error) {
	if Full(buf) {
		buf.WrapReader.Seek(buf.WrapWriter.Offset(), 0)
	}
	return buf.WrapReader.ReadAt(p, buf.WrapReader.Offset()+off)
}

func (buf *Ring) Write(p []byte) (n int, err error) {
	n, err = buf.WrapWriter.Write(p)
	buf.L += int64(n)
	if buf.L > buf.Cap() {
		buf.L = buf.Cap()
	}
	return n, err
}

func (buf *Ring) WriteAt(p []byte, off int64) (n int, err error) {
	n, err = buf.WrapWriter.WriteAt(p, off)
	if off+int64(n) > buf.Len() {
		buf.L = off + int64(n)
		if buf.L > buf.Cap() {
			buf.L = buf.Cap()
		}
	}
	return n, err
}

func (buf *Ring) FFwd(off int) int {
	if Full(buf) {
		buf.WrapReader.Seek(buf.WrapWriter.Offset(), 0)
	}
	if int64(off) >= buf.Len() {
		buf.Reset()
		return int(buf.Len())
	}
	buf.WrapReader.Seek(int64(off), 1)
	buf.L -= int64(off)
	return off
}

func (buf *Ring) Reset() {
	buf.Buffer.Reset()
	buf.L = 0
	buf.WrapReader = NewWrapReader(buf.Buffer, 0, buf.Cap())
	buf.WrapWriter = NewWrapWriter(buf.Buffer, 0, buf.Cap())
}
