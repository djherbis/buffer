package buffer

import (
	"io"
)

const (
	maxUint   = ^uint32(0)
	maxInt    = int32(maxUint >> 1)
	maxUint64 = ^uint64(0)
	maxInt64  = int64(maxUint64 >> 1)
	maxAlloc  = int64(1024 * 32)
)

type Resetter interface {
	Reset()
}

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
	Lener
	Caper
	Resetter
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
	FFwd(int64) int64
}

type buffer struct {
	head []byte
	Buffer
}

func len64(p []byte) int64 {
	return int64(len(p))
}

func MaxCap() int64 {
	return maxInt64
}

func LimitAlloc(n int64) int64 {
	if n > maxAlloc {
		return maxAlloc
	}
	return n
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

func RoomFor(buf LenCaper, p []byte) bool {
	return Gap(buf) >= len64(p)
}

func ShrinkToRead(buf Lener, p []byte) []byte {
	if buf.Len() < len64(p) {
		return p[:buf.Len()]
	}
	return p
}

func ShrinkToFit(buf LenCaper, p []byte) []byte {
	if !RoomFor(buf, p) {
		p = p[:Gap(buf)]
	}
	return p
}

func TotalLen(buffers []Buffer) (n int64) {
	for _, buffer := range buffers {
		if n > MaxCap()-buffer.Len() {
			return MaxCap()
		} else {
			n += buffer.Len()
		}
	}
	return n
}

func TotalCap(buffers []Buffer) (n int64) {
	for _, buffer := range buffers {
		if n > MaxCap()-buffer.Cap() {
			return MaxCap()
		} else {
			n += buffer.Cap()
		}
	}
	return n
}

func ResetAll(buffers []Buffer) {
	for _, buffer := range buffers {
		buffer.Reset()
	}
}

func NewUnboundedBuffer(mem, chunk int64) Buffer {
	return NewMulti(New(mem),
		NewPartition(func() Buffer { return NewFile(chunk) }))
}
