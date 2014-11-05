package buffer

import (
	"bytes"
	"io"
)

type LinkBuffer struct {
	buffer  Buffer
	next    Buffer
	hasNext bool
}

func NewMulti(buffers ...Buffer) *LinkBuffer {
	if len(buffers) == 0 {
		return nil
	}

	buf := &LinkBuffer{
		buffer:  buffers[0],
		next:    NewMulti(buffers[1:]...),
		hasNext: len(buffers[1:]) != 0,
	}

	return buf
}

func (buf *LinkBuffer) Reset() {
	if buf.hasNext {
		buf.next.Reset()
	}
	buf.buffer.Reset()
}

func (buf *LinkBuffer) Cap() (n int64) {
	if buf.hasNext {
		next := buf.next.Cap()
		if buf.buffer.Cap() > MaxCap()-next {
			return MaxCap()
		} else {
			return buf.buffer.Cap() + next
		}
	}

	return buf.buffer.Cap()
}

func (buf *LinkBuffer) Len() (n int64) {
	if buf.hasNext {
		next := buf.next.Len()
		if buf.buffer.Len() > MaxCap()-next {
			return MaxCap()
		} else {
			return buf.buffer.Len() + next
		}
	}

	return buf.buffer.Len()
}

func (buf *LinkBuffer) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, io.EOF
}

func (buf *LinkBuffer) Read(p []byte) (n int, err error) {
	n, err = buf.buffer.Read(p)
	p = p[n:]
	if len(p) > 0 && buf.hasNext && (err == nil || err == io.EOF) {
		m, err := buf.next.Read(p)
		n += m
		if err != nil {
			return n, err
		}
	}

	for !Full(buf.buffer) && buf.hasNext && !Empty(buf.next) {
		r := io.LimitReader(buf.next, Gap(buf.buffer))
		if _, err = io.Copy(buf.buffer, r); err != nil && err != io.EOF {
			return n, err
		}
	}
	return n, err
}

func (buf *LinkBuffer) Write(p []byte) (n int, err error) {
	if n, err = buf.buffer.Write(p); err == bytes.ErrTooLarge {
		err = nil
	}
	p = p[n:]
	if len(p) > 0 && buf.hasNext && err == nil {
		m, err := buf.next.Write(p)
		n += m
		if err != nil {
			return n, err
		}
	}
	return n, err
}
