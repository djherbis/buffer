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

type Buffer interface {
	Len() int64
	Cap() int64
	io.Reader
	io.Writer
}

type buffer struct {
	head []byte
	Buffer
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

func Gap(buf Buffer) int64 {
	return buf.Cap() - buf.Len()
}

func Full(buf Buffer) bool {
	return Gap(buf) == 0
}

func Empty(buf Buffer) bool {
	return buf.Len() == 0
}

func RoomFor(buf Buffer, p []byte) bool {
	return Gap(buf) >= int64(len(p))
}

func ShrinkToRead(buf Buffer, p []byte) []byte {
	if buf.Len() < int64(len(p)) {
		return p[:buf.Len()]
	}
	return p
}

func ShrinkToFit(buf Buffer, p []byte) []byte {
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

func NewUnboundedBuffer(mem, chunk int64) Buffer {
	return NewMulti(New(mem), NewPartition(chunk, NewFile))
}
