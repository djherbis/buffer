package buffer

import (
	"bytes"
	"encoding/gob"
	"io"
)

type chain struct {
	Buf     Buffer
	Next    Buffer
	HasNext bool
}

// NewMulti returns a Buffer which is the logical concatenation of the passed buffers.
func NewMulti(buffers ...Buffer) Buffer {
	if len(buffers) == 0 {
		return nil
	} else if len(buffers) == 1 {
		return buffers[0]
	}

	buf := &chain{
		Buf:     buffers[0],
		Next:    NewMulti(buffers[1:]...),
		HasNext: len(buffers[1:]) != 0,
	}

	buf.Defrag()

	return buf
}

func (buf *chain) Reset() {
	if buf.HasNext {
		buf.Next.Reset()
	}
	buf.Buf.Reset()
}

func (buf *chain) Cap() (n int64) {
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

func (buf *chain) Len() (n int64) {
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

func (buf *chain) Defrag() {
	for !Full(buf.Buf) && buf.HasNext && !Empty(buf.Next) {
		r := io.LimitReader(buf.Next, Gap(buf.Buf))
		if _, err := io.Copy(buf.Buf, r); err != nil && err != io.EOF {
			return
		}
	}
}

func (buf *chain) Read(p []byte) (n int, err error) {
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

func (buf *chain) Write(p []byte) (n int, err error) {
	if n, err = buf.Buf.Write(p); err == io.ErrShortWrite && buf.HasNext {
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
	gob.Register(&chain{})
}

func (buf *chain) MarshalBinary() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(b)
	enc.Encode(&buf.Buf)
	enc.Encode(buf.HasNext)
	if buf.HasNext {
		enc.Encode(&buf.Next)
	}
	return b.Bytes(), nil
}

func (buf *chain) UnmarshalBinary(data []byte) error {
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
