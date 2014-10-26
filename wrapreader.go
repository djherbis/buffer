package buffer

import "io"

type wrapReader struct {
	off    int64
	wrapAt int64
	io.ReadSeeker
}

func NewWrapReader(w io.ReadSeeker, wrapAt int64) io.Reader {
	return &wrapReader{
		ReadSeeker: w,
		wrapAt:     wrapAt,
	}
}

func (w *wrapReader) Read(p []byte) (n int, err error) {
	w.Seek(w.off, 0)
	n, err = w.ReadSeeker.Read(p)
	w.off += int64(n)
	return n, err
}
