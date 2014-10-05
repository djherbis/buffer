package buffer

import (
	"io"
)

type partition struct {
	buffers []Buffer
	make    func(int64) Buffer
	chunk   int64
}

func NewPartition(chunk int64, make func(int64) Buffer) Buffer {
	buf := &partition{
		make:  make,
		chunk: chunk,
	}
	buf.push()
	return newBuffer(buf)
}

func (buf *partition) Len() int64 {
	return TotalLen(buf.buffers)
}

func (buf *partition) Cap() int64 {
	return MaxCap()
}

func (buf *partition) Read(p []byte) (n int, err error) {
	for len(p) > 0 {

		if len(buf.buffers) == 0 {
			return n, io.EOF
		}

		buffer := buf.buffers[0]

		if Empty(buffer) {
			buf.pop()
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

		if len(buf.buffers) == 0 {
			buf.push()
		}

		buffer := buf.buffers[len(buf.buffers)-1]

		if Full(buffer) {
			buf.push()
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

func (buf *partition) push() {
	buf.buffers = append(buf.buffers, buf.make(buf.chunk))
}

func (buf *partition) pop() {
	buf.buffers = buf.buffers[1:]
}
