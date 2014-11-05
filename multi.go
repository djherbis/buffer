package buffer

import (
	"bytes"
	"encoding/gob"
	"io"
)

type LinkBuffer struct {
	Buf     Buffer
	Next    Buffer
	HasNext bool
}

func NewMulti(buffers ...Buffer) *LinkBuffer {
	if len(buffers) == 0 {
		return nil
	}

	buf := &LinkBuffer{
		Buf:     buffers[0],
		Next:    NewMulti(buffers[1:]...),
		HasNext: len(buffers[1:]) != 0,
	}

	return buf
}

func (buf *LinkBuffer) Reset() {
	if buf.HasNext {
		buf.Next.Reset()
	}
	buf.Buf.Reset()
}

func (buf *LinkBuffer) Cap() (n int64) {
	if buf.HasNext {
		Next := buf.Next.Cap()
		if buf.Buf.Cap() > MaxCap()-Next {
			return MaxCap()
		} else {
			return buf.Buf.Cap() + Next
		}
	}

	return buf.Buf.Cap()
}

func (buf *LinkBuffer) Len() (n int64) {
	if buf.HasNext {
		Next := buf.Next.Len()
		if buf.Buf.Len() > MaxCap()-Next {
			return MaxCap()
		} else {
			return buf.Buf.Len() + Next
		}
	}

	return buf.Buf.Len()
}

func (buf *LinkBuffer) ReadAt(p []byte, off int64) (n int, err error) {
	n, err = buf.Buf.ReadAt(p, off)
	p = p[n:]

	if len(p) > 0 && buf.HasNext && (err == nil || err == io.EOF) {
		if off > buf.Buf.Len() {
			off -= buf.Buf.Len()
		} else {
			off = 0
		}
		m, err := buf.Next.ReadAt(p, off)
		n += m
		if err != nil {
			return n, err
		}
	}

	return n, err
}

func (buf *LinkBuffer) FFwd(n int) int {
	m := buf.Buf.FFwd(n)

	if n > m && buf.HasNext {
		m += buf.Next.FFwd(n - m)
	}

	for !Full(buf.Buf) && buf.HasNext && !Empty(buf.Next) {
		r := io.LimitReader(buf.Next, Gap(buf.Buf))
		if _, err := io.Copy(buf.Buf, r); err != nil && err != io.EOF {
			return m
		}
	}

	return m
}

func (buf *LinkBuffer) Read(p []byte) (n int, err error) {
	n, err = buf.Buf.Read(p)
	p = p[n:]
	if len(p) > 0 && buf.HasNext && (err == nil || err == io.EOF) {
		m, err := buf.Next.Read(p)
		n += m
		if err != nil {
			return n, err
		}
	}

	for !Full(buf.Buf) && buf.HasNext && !Empty(buf.Next) {
		r := io.LimitReader(buf.Next, Gap(buf.Buf))
		if _, err = io.Copy(buf.Buf, r); err != nil && err != io.EOF {
			return n, err
		}
	}
	return n, err
}

func (buf *LinkBuffer) Write(p []byte) (n int, err error) {
	if n, err = buf.Buf.Write(p); err == bytes.ErrTooLarge && buf.HasNext {
		err = nil
	}
	p = p[n:]
	if len(p) > 0 && buf.HasNext && err == nil {
		m, err := buf.Next.Write(p)
		n += m
		if err != nil {
			return n, err
		}
	}
	return n, err
}

func init() {
	gob.Register(&LinkBuffer{})
}

func (buf *LinkBuffer) MarshalBinary() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(b)
	enc.Encode(&buf.Buf)
	enc.Encode(buf.HasNext)
	if buf.HasNext {
		enc.Encode(&buf.Next)
	}
	return b.Bytes(), nil
}

func (buf *LinkBuffer) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	if err := dec.Decode(&buf.Buf); err != nil {
		return err
	}
	if err := dec.Decode(&buf.HasNext); err != nil {
		return err
	}
	if buf.HasNext {
		if err := dec.Decode(&buf.Next); err != nil {
			return err
		}
	}
	return nil
}
