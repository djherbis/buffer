package buffer

import (
	"bytes"
	"encoding/gob"
	"io"
)

type Chain struct {
	Buf     Buffer
	Next    Buffer
	HasNext bool
}

func NewMulti(buffers ...Buffer) Buffer {
	if len(buffers) == 0 {
		return nil
	} else if len(buffers) == 1 {
		return buffers[0]
	}

	buf := &Chain{
		Buf:     buffers[0],
		Next:    NewMulti(buffers[1:]...),
		HasNext: len(buffers[1:]) != 0,
	}

	buf.Defrag()

	return buf
}

func (buf *Chain) Reset() {
	if buf.HasNext {
		buf.Next.Reset()
	}
	buf.Buf.Reset()
}

func (buf *Chain) Cap() (n int64) {
	if buf.HasNext {
		Next := buf.Next.Cap()
		if buf.Buf.Cap() > MAXINT64-Next {
			return MAXINT64
		} else {
			return buf.Buf.Cap() + Next
		}
	}

	return buf.Buf.Cap()
}

func (buf *Chain) Len() (n int64) {
	if buf.HasNext {
		Next := buf.Next.Len()
		if buf.Buf.Len() > MAXINT64-Next {
			return MAXINT64
		} else {
			return buf.Buf.Len() + Next
		}
	}

	return buf.Buf.Len()
}

func (buf *Chain) ReadAt(p []byte, off int64) (n int, err error) {
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

func (buf *Chain) Defrag() {
	for !Full(buf.Buf) && buf.HasNext && !Empty(buf.Next) {
		r := io.LimitReader(buf.Next, Gap(buf.Buf))
		if _, err := io.Copy(buf.Buf, r); err != nil && err != io.EOF {
			return
		}
	}
}

func (buf *Chain) Read(p []byte) (n int, err error) {
	n, err = buf.Buf.Read(p)
	p = p[n:]
	if len(p) > 0 && buf.HasNext && (err == nil || err == io.EOF) {
		m, err := buf.Next.Read(p)
		n += m
		if err != nil {
			return n, err
		}
	}

	buf.Defrag()

	return n, err
}

func (buf *Chain) Write(p []byte) (n int, err error) {
	if n, err = buf.Buf.Write(p); err == io.ErrShortBuffer && buf.HasNext {
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

func (buf *Chain) WriteAt(p []byte, off int64) (n int, err error) {
	n, err = buf.Buf.WriteAt(p, off)
	p = p[n:]

	if len(p) > 0 && buf.HasNext && (err == nil || err == io.ErrShortBuffer) {
		if off > buf.Buf.Len() {
			off -= buf.Buf.Len()
		} else {
			off = 0
		}

		var m int
		m, err = buf.Next.WriteAt(p, off)
		n += m
		if err != nil {
			return n, err
		}
	}

	return n, err
}

func init() {
	gob.Register(&Chain{})
}

func (buf *Chain) MarshalBinary() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(b)
	enc.Encode(&buf.Buf)
	enc.Encode(buf.HasNext)
	if buf.HasNext {
		enc.Encode(&buf.Next)
	}
	return b.Bytes(), nil
}

func (buf *Chain) UnmarshalBinary(data []byte) error {
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
