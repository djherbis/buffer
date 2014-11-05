package buffer

import "io"

type DoerAt interface {
	DoAt([]byte, int64) (int, error)
}
type DoAtFunc func([]byte, int64) (int, error)

type wrapper struct {
	off    int64
	wrapAt int64
	doat   DoAtFunc
}

func (w *wrapper) DoAt(p []byte, off int64) (n int, err error) {
	return w.doat(p, off)
}

type WrapWriter struct {
	*wrapper
}

func NewWrapWriter(w io.WriterAt, off int64, wrapAt int64) *WrapWriter {
	return &WrapWriter{
		&wrapper{
			doat:   w.WriteAt,
			off:    off,
			wrapAt: wrapAt,
		},
	}
}

func (w *WrapWriter) Write(p []byte) (n int, err error) {
	n, err = Wrap(w, p, w.off, w.wrapAt)
	w.off = (w.off + int64(n)) % w.wrapAt
	return n, err
}

func (w *WrapWriter) WriteAt(p []byte, off int64) (n int, err error) {
	return Wrap(w, p, off, w.wrapAt)
}

type WrapReader struct {
	*wrapper
}

func NewWrapReader(r io.ReaderAt, off int64, wrapAt int64) *WrapReader {
	return &WrapReader{
		&wrapper{
			doat:   r.ReadAt,
			off:    off,
			wrapAt: wrapAt,
		},
	}
}

func (r *WrapReader) Read(p []byte) (n int, err error) {
	n, err = Wrap(r, p, r.off, r.wrapAt)
	r.off = (r.off + int64(n)) % r.wrapAt
	return n, err
}

func (r *WrapReader) ReadAt(p []byte, off int64) (n int, err error) {
	return Wrap(r, p, off, r.wrapAt)
}

func Wrap(w DoerAt, p []byte, off int64, wrapAt int64) (n int, err error) {
	var m int

	for len(p) > 0 {

		if off+len64(p) <= wrapAt {
			if m, err = w.DoAt(p, off); err != nil {
				return n + m, err
			}
		} else {
			space := wrapAt - off
			if m, err = w.DoAt(p[:space], off); err != nil {
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
