package buffer

import (
	"io"
)

type linkBuffer struct {
	buffer  Buffer
	next    Buffer
	hasNext bool
}

func NewMulti(buffers ...Buffer) Buffer {
	if len(buffers) == 0 {
		return nil
	}

	buf := &linkBuffer{
		buffer:  buffers[0],
		next:    NewMulti(buffers[1:]...),
		hasNext: len(buffers[1:]) != 0,
	}

	return newBuffer(buf)
}

func (buf *linkBuffer) Cap() (n int64) {
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

func (buf *linkBuffer) Len() (n int64) {
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

func (buf *linkBuffer) Read(p []byte) (n int, err error) {
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
		gap := Gap(buf.buffer)
		r := io.LimitReader(buf.next, gap)
		if _, err = io.Copy(buf.buffer, r); err != nil && err != io.EOF {
			return n, err
		}
	}
	return n, err
}

func (buf *linkBuffer) Write(p []byte) (n int, err error) {
	n, err = buf.buffer.Write(p)
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
