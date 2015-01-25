package buffer

import (
	"encoding/gob"
	"io"
)

type partition struct {
	BufferList
	Pool
}

// NewPartition returns a Buffer which uses a Pool to extend-shrink its size as needed.
// It automatically allocates new buffers with pool.Get() to extend is length, and
// pool.Put() to release unused buffers as it shrinks.
func NewPartition(pool Pool, buffers ...Buffer) Buffer {
	return &partition{
		Pool:       pool,
		BufferList: buffers,
	}
}

func (buf *partition) Cap() int64 {
	return MAXINT64
}

func (buf *partition) Read(p []byte) (n int, err error) {
	for len(p) > 0 {

		if len(buf.BufferList) == 0 {
			return n, io.EOF
		}

		buffer := buf.BufferList[0]

		if Empty(buffer) {
			buf.Pool.Put(buf.Pop())
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
			buf.Push(buf.Pool.Get())
		}

		buffer := buf.BufferList[len(buf.BufferList)-1]

		if Full(buffer) {
			buf.Push(buf.Pool.Get())
			continue
		}

		m, er := buffer.Write(p)
		n += m
		p = p[m:]

		if er == io.ErrShortBuffer {
			er = nil
		} else if er != nil {
			return n, er
		}

	}
	return n, nil
}

func (buf *partition) Reset() {
	for len(buf.BufferList) > 0 {
		buf.Pool.Put(buf.Pop())
	}
}

func init() {
	gob.Register(&partition{})
}
