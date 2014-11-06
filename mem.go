package buffer

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

type MemBuffer struct {
	N int64
	*bytes.Buffer
}

func New(n int64) *MemBuffer {
	return &MemBuffer{
		N:      n,
		Buffer: bytes.NewBuffer(nil),
	}
}

func (buf *MemBuffer) Cap() int64 {
	return buf.N
}

func (buf *MemBuffer) Len() int64 {
	return int64(buf.Buffer.Len())
}

func (buf *MemBuffer) Write(p []byte) (n int, err error) {
	return LimitWriter(buf.Buffer, Gap(buf)).Write(p)
}

// Must Add io.ErrShortWrite
func (buf *MemBuffer) WriteAt(p []byte, off int64) (n int, err error) {
	if off > buf.Len() {
		return 0, bytes.ErrTooLarge
	} else if len64(p)+off <= buf.Len() {
		d := buf.Bytes()[off:]
		return copy(d, p), nil
	} else {
		d := buf.Bytes()[off:]
		n = copy(d, p)
		m, err := buf.Write(p[n:])
		return n + m, err
	}
}

func (buf *MemBuffer) ReadAt(p []byte, off int64) (n int, err error) {
	return bytes.NewReader(buf.Bytes()).ReadAt(p, off)
}

func (buf *MemBuffer) Read(p []byte) (n int, err error) {
	return io.LimitReader(buf.Buffer, buf.Len()).Read(p)
}

func (buf *MemBuffer) FFwd(n int) int {
	data := buf.Bytes()
	if n > len(data) {
		n = len(data)
	}
	b := bytes.NewBuffer(data[n:])
	buf.Buffer = b
	return n
}

func init() {
	gob.Register(&MemBuffer{})
}

func (buf *MemBuffer) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprint(&b, buf.N)
	b.Write(buf.Bytes())
	return b.Bytes(), nil
}

func (buf *MemBuffer) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscan(b, &buf.N)
	buf.Buffer = b
	return err
}
