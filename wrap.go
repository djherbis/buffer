package buffer

import "io"

type SeekFunc func(int64, int) (int64, error)
type ByteFunc func([]byte) (int, error)

type wrapper struct {
	off    int64
	wrapAt int64
	seek   SeekFunc
	do     ByteFunc
}

type WrapWriter struct {
	wrapper
}

func NewWrapWriter(w io.WriteSeeker, wrapAt int64) *WrapWriter {
	return &WrapWriter{wrapper{
		seek:   w.Seek,
		do:     w.Write,
		wrapAt: wrapAt,
	}}
}

func (w *WrapWriter) Write(p []byte) (n int, err error) {
	n, err = w.Wrap(p, w.off)
	w.off = (w.off + int64(n)) % w.wrapAt
	return n, err
}

func (w *WrapWriter) WriteAt(p []byte, off int64) (n int, err error) {
	return w.Wrap(p, off)
}

type WrapReader struct {
	wrapper
}

func NewWrapReader(r io.ReadSeeker, wrapAt int64) *WrapReader {
	return &WrapReader{wrapper{
		seek:   r.Seek,
		do:     r.Read,
		wrapAt: wrapAt,
	}}
}

func (w *WrapReader) Read(p []byte) (n int, err error) {
	n, err = w.Wrap(p, w.off)
	w.off = (w.off + int64(n)) % w.wrapAt
	return n, err
}

func (w *WrapReader) ReadAt(p []byte, off int64) (n int, err error) {
	return w.Wrap(p, off)
}

func (w *wrapper) Seek(off int64, rel int) (int64, error) {
	return w.seek(off, rel)
}

func (w *wrapper) Wrap(p []byte, off int64) (n int, err error) {
	var m int

	for len(p) > 0 {

		w.seek(off, 0)

		if off+len64(p) <= w.wrapAt {
			if m, err = w.do(p); err != nil {
				return n + m, err
			}
		} else {
			space := w.wrapAt - off
			if m, err = w.do(p[:space]); err != nil {
				return n + m, err
			}
		}

		n += m
		off += int64(m)
		p = p[m:]

		if off == w.wrapAt {
			off = 0
		}
	}

	return n, err
}
