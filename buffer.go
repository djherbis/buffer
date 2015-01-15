package buffer

import (
	"io"
)

const MAXINT64 = 9223372036854775807

type Buffer interface {
	Len() int64
	Cap() int64
	io.Reader
	io.Writer
	Reset()
}

type BufferAt interface {
	Buffer
	io.ReaderAt
	io.WriterAt
}

func len64(p []byte) int64 {
	return int64(len(p))
}

func Gap(buf Buffer) int64 {
	return buf.Cap() - buf.Len()
}

func Full(buf Buffer) bool {
	return Gap(buf) == 0
}

func Empty(l Buffer) bool {
	return l.Len() == 0
}

func NewUnboundedBuffer(mem, chunk int64) Buffer {
	return NewMulti(New(mem),
		NewPartition(func() Buffer { return NewFile(chunk) }))
}
