package buffer

import "io"

type LimitedWriter struct {
	W io.Writer
	N int64
}

func (l *LimitedWriter) Write(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, io.ErrShortBuffer
	}
	if len64(p) > l.N {
		p = p[0:l.N]
	}
	n, err = l.W.Write(p)
	l.N -= int64(n)
	return n, err
}

func LimitWriter(w io.Writer, n int64) io.Writer {
	return &LimitedWriter{W: w, N: n}
}
