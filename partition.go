package buffer

import (
	"bytes"
	"encoding/gob"
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
	return MAXINT64
}

func (buf *Partition) ReadAt(p []byte, off int64) (n int, err error) {
	index := 0
	for off > 0 && index < len(buf.BufferList) {
		buffer := buf.BufferList[index]
		if off >= buffer.Len() {
			off -= buffer.Len()
			index++
		} else {
			break
		}
	}

	for len(p) > 0 {

		if index >= len(buf.BufferList) {
			return n, io.EOF
		}

		buffer := buf.BufferList[index]

		m, er := buffer.ReadAt(p, off)
		n += m
		p = p[m:]

		if er == io.EOF || int64(m)+off == buffer.Len() {
			index++
			off = 0
		} else if er != nil {
			return n, er
		}

	}

	return n, nil
}

func (buf *Partition) FFwd(n int64) (m int64) {
	for n > m {
		if len(buf.BufferList) == 0 {
			return m
		}

		buffer := buf.BufferList[0]

		if Empty(buffer) {
			buf.Pop()
			continue
		}

		m += buffer.FFwd(n - m)
	}

	return m
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

		if er == bytes.ErrTooLarge {
			er = nil
		} else if er != nil {
			return n, er
		}

	}
	return n, nil
}

func (buf *Partition) WriteAt(p []byte, off int64) (n int, err error) {
	index := 0
	for off > 0 && index < len(buf.BufferList) {
		buffer := buf.BufferList[index]
		if off >= buffer.Len() {
			off -= buffer.Len()
			index++
		} else {
			break
		}
	}

	for len(p) > 0 {

		if index >= len(buf.BufferList) {
			if off > 0 {
				return n, bytes.ErrTooLarge
			}
			buf.BufferList.Push(buf.make())
		}

		buffer := buf.BufferList[index]

		m, er := buffer.WriteAt(p, off)
		n += m
		p = p[m:]

		if er == bytes.ErrTooLarge || int64(m)+off == buffer.Cap() {
			index++
			off = 0
		} else if er != nil {
			return n, er
		}

	}

	return n, nil
}

func init() {
	gob.Register(&Partition{})
}
