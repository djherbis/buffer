package buffer

import "io"

type wrapWriter struct {
	off    int64
	wrapAt int64
	io.WriteSeeker
}

func NewWrapWriter(w io.WriteSeeker, wrapAt int64) io.Writer {
	return &wrapWriter{
		WriteSeeker: w,
		wrapAt:      wrapAt,
	}
}

func (w *wrapWriter) Write(p []byte) (n int, err error) {
	w.Seek(w.off, 0)
	n, err = w.WriteSeeker.Write(p)
	w.off += int64(n)
	return n, err
}
