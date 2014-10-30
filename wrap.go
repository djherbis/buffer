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

type wrapWriter struct {
	wrapper
}

func NewWrapWriter(w io.WriteSeeker, wrapAt int64) io.Writer {
	return &wrapWriter{wrapper{
		seek:   w.Seek,
		do:     w.Write,
		wrapAt: wrapAt,
	}}
}

func (w *wrapWriter) Write(p []byte) (n int, err error) {
	return w.Wrap(p)
}

type wrapReader struct {
	wrapper
}

func NewWrapReader(r io.ReadSeeker, wrapAt int64) io.Reader {
	return &wrapReader{wrapper{
		seek:   r.Seek,
		do:     r.Read,
		wrapAt: wrapAt,
	}}
}

func (w *wrapReader) Read(p []byte) (n int, err error) {
	return w.Wrap(p)
}

func (w *wrapper) Wrap(p []byte) (n int, err error) {
	var m int

	for len(p) > 0 {

		w.seek(w.off, 0)

		if w.off+len64(p) <= w.wrapAt {
			if m, err = w.do(p); err != nil {
				return n + m, err
			}
		} else {
			space := w.wrapAt - w.off
			if m, err = w.do(p[:space]); err != nil {
				return n + m, err
			}
		}

		n += m
		w.off += int64(m)
		p = p[m:]

		if w.off == w.wrapAt {
			w.off = 0
		}
	}

	return n, err
}
