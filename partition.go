package buffer

import (
	"io"
)

type partition struct {
	BufferList
	make  func(int64) Buffer
	chunk int64
}

func NewPartition(chunk int64, make func(int64) Buffer) Buffer {
	buf := &partition{
		make:  make,
		chunk: chunk,
	}
	buf.Push(buf.make(buf.chunk))
	return buf
}

func (buf *partition) Cap() int64 {
	return MaxCap()
}

func (buf *partition) Read(p []byte) (n int, err error) {
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

func (buf *partition) Write(p []byte) (n int, err error) {
	for len(p) > 0 {

		if len(buf.BufferList) == 0 {
			buf.Push(buf.make(buf.chunk))
		}

		buffer := buf.BufferList[len(buf.BufferList)-1]

		if Full(buffer) {
			buf.Push(buf.make(buf.chunk))
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
