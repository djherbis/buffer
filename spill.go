package buffer

import (
	"encoding/gob"
	"io"
)

type Spill struct {
	Buffer
	Spiller io.Writer
}

func NewSpill(buf Buffer, w io.Writer) Buffer {
	return &Spill{
		Buffer:  buf,
		Spiller: w,
	}
}

func (buf *Spill) Cap() int64 {
	return MAXINT64
}

func (buf *Spill) Write(p []byte) (n int, err error) {
	if n, err = buf.Buffer.Write(p); err != nil {
		// BUG(Dustin): What is this fails? What if it doesn't write everything?
		buf.Spiller.Write(p[:n])
	}
	return len(p), nil
}

func init() {
	gob.Register(&Spill{})
}
