package buffer

import "io"

type SpillBuffer struct {
	Buffer
	Spiller io.Writer
}

func (buf *SpillBuffer) Write(p []byte) (n int, err error) {
	if n, err = buf.Buffer.Write(p); err != nil {
		buf.Spiller.Write(p[:n])
	}
	return len(p), nil
}
