package buffer

import (
	"io"
)

const MAXINT64 = 9223372036854775807

type Lener interface {
	Len() int64
}

type Caper interface {
	Cap() int64
}

type LenCaper interface {
	Lener
	Caper
}

type Buffer interface {
	LenCaper
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
	Reset()
}

func len64(p []byte) int64 {
	return int64(len(p))
}

func Gap(buf LenCaper) int64 {
	return buf.Cap() - buf.Len()
}

func Full(buf LenCaper) bool {
	return Gap(buf) == 0
}

func Empty(l Lener) bool {
	return l.Len() == 0
}

func NewUnboundedBuffer(mem, chunk int64) Buffer {
	return NewMulti(New(mem),
		NewPartition(func() Buffer { return NewFile(chunk) }))
}
