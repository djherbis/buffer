package buffer

import (
	"encoding/gob"
	"io"
)

type spill struct {
	Buffer
	Spiller io.Writer
}

func NewSpill(buf Buffer, w io.Writer) Buffer {
	return &spill{
		Buffer:  buf,
		Spiller: w,
	}
}

func (buf *spill) Cap() int64 {
	return MAXINT64
}

func (buf *spill) Write(p []byte) (n int, err error) {
	if n, err = buf.Buffer.Write(p); err != nil {
		// BUG(Dustin): What is this fails? What if it doesn't write everything?
		buf.Spiller.Write(p[:n])
	}
	return len(p), nil
}

func init() {
	gob.Register(&spill{})
}
