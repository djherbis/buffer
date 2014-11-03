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

type WrapWriter wrapper

func NewWrapWriter(w io.WriteSeeker, wrapAt int64) *WrapWriter {
	return &WrapWriter{
		seek:   w.Seek,
		do:     w.Write,
		wrapAt: wrapAt,
	}
}

func (w *WrapWriter) Write(p []byte) (n int, err error) {
	n, err = Wrap(w.seek, w.do, p, w.off, w.wrapAt)
	w.off = (w.off + int64(n)) % w.wrapAt
	return n, err
}

func (w *WrapWriter) WriteAt(p []byte, off int64) (n int, err error) {
	return Wrap(w.seek, w.do, p, off, w.wrapAt)
}

type WrapReader wrapper

func NewWrapReader(r io.ReadSeeker, wrapAt int64) *WrapReader {
	return &WrapReader{
		seek:   r.Seek,
		do:     r.Read,
		wrapAt: wrapAt,
	}
}

func (r *WrapReader) Read(p []byte) (n int, err error) {
	n, err = Wrap(r.seek, r.do, p, r.off, r.wrapAt)
	r.off = (r.off + int64(n)) % r.wrapAt
	return n, err
}

func (r *WrapReader) ReadAt(p []byte, off int64) (n int, err error) {
	return Wrap(r.seek, r.do, p, off, r.wrapAt)
}

func (w *wrapper) Seek(off int64, rel int) (int64, error) {
	return w.seek(off, rel)
}

func Wrap(seek SeekFunc, do ByteFunc, p []byte, off int64, wrapAt int64) (n int, err error) {
	var m int

	for len(p) > 0 {

		seek(off, 0)

		if off+len64(p) <= wrapAt {
			if m, err = do(p); err != nil {
				return n + m, err
			}
		} else {
			space := wrapAt - off
			if m, err = do(p[:space]); err != nil {
				return n + m, err
			}
		}

		n += m
		off += int64(m)
		p = p[m:]

		if off == wrapAt {
			off = 0
		}
	}

	return n, err
}
