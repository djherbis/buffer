package buffer

import (
	"io"
)

type Partition struct {
	BufferList
	make func() Buffer
}

func NewPartition(make func() Buffer, buffers ...Buffer) *Partition {
	buf := &Partition{
		make:       make,
		BufferList: buffers,
	}
	return buf
}

func (buf *Partition) Cap() int64 {
	return MaxCap()
}

func (buf *Partition) Read(p []byte) (n int, err error) {
	for len(p) > 0 {

		if len(buf.BufferList) == 0 {
			return n, io.EOF
		}

		buffer := buf.BufferList[0]

		if Empty(buffer) {
			buf.Pop()
			continue
		}

		m, er := buffer.Read(p)
		n += m
		p = p[m:]

		if er != nil && er != io.EOF {
			return n, er
		}

	}
	return n, nil
}

func (buf *Partition) Write(p []byte) (n int, err error) {
	for len(p) > 0 {

		if len(buf.BufferList) == 0 {
			buf.Push(buf.make())
		}

		buffer := buf.BufferList[len(buf.BufferList)-1]

		if Full(buffer) {
			buf.Push(buf.make())
			continue
		}

		m, er := buffer.Write(p)
		n += m
		p = p[m:]

		if er != nil {
			return n, er
		}

	}
	return n, nil
}
