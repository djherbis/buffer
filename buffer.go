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

func newBuffer(buf Buffer) Buffer {
	return &buffer{
		Buffer: buf,
	}
}

func (buf *buffer) Len() int64 {
	return buf.Buffer.Len() + int64(len(buf.head))
}

func (buf *buffer) Write(p []byte) (n int, err error) {
	n, err = buf.Buffer.Write(ShrinkToFit(buf, p))
	return n, err
}

func (buf *buffer) WriteTo(w io.Writer) (n int64, err error) {
	for !Empty(buf) || len(buf.head) != 0 {
		if len(buf.head) == 0 {
			buf.head = make([]byte, 1024*32)
			m, er := buf.Buffer.Read(buf.head)
			buf.head = buf.head[:m]
			if er != nil && er != io.EOF {
				return n, er
			}
		}

		m, er := w.Write(buf.head)
		n += int64(m)
		buf.head = buf.head[m:]
		if er != nil {
			return n, er
		}

	}
	return n, nil
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