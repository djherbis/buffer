package wrapio

import (
	"encoding/gob"
	"io"

	"github.com/djherbis/buffer/limio"
)

type ReadWriterAt interface {
	io.ReaderAt
	io.WriterAt
}

type Wrapper struct {
	N   int64
	L   int64
	O   int64
	rwa ReadWriterAt
}

func NewWrapper(rwa ReadWriterAt, N int64) *Wrapper {
	return &Wrapper{
		N:   N,
		rwa: rwa,
	}
}

func (wpr *Wrapper) Len() int64 {
	return wpr.L
}

func (wpr *Wrapper) Cap() int64 {
	return wpr.N
}

func (wpr *Wrapper) Reset() {
	wpr.O = 0
	wpr.L = 0
}

func (wpr *Wrapper) Read(p []byte) (n int, err error) {
	n, err = wpr.ReadAt(p, 0)
	wpr.L -= int64(n)
	wpr.O += int64(n)
	wpr.O %= wpr.N
	return n, err
}

// BUG(Dustin): Handle if off is too far / len(p) too big
func (wpr *Wrapper) ReadAt(p []byte, off int64) (n int, err error) {
	wrap := NewWrapReader(wpr.rwa, wpr.O+off, wpr.N)
	r := io.LimitReader(wrap, wpr.L-off)
	return r.Read(p)
}

func (wpr *Wrapper) Write(p []byte) (n int, err error) {
	return wpr.WriteAt(p, wpr.L)
}

// BUG(Dustin): Handle if off is too far / len(p) too big
func (wpr *Wrapper) WriteAt(p []byte, off int64) (n int, err error) {
	wrap := NewWrapWriter(wpr.rwa, wpr.O+off, wpr.N)
	w := limio.LimitWriter(wrap, wpr.N-off)
	n, err = w.Write(p)
	if wpr.L < off+int64(n) {
		wpr.L = int64(n) + off
	}
	return n, err
}

func init() {
	gob.Register(&Wrapper{})
}
