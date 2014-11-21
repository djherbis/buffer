package buffer

import (
	"encoding/gob"
	"io"
)

type Spill struct {
	Buffer
	Spiller io.Writer
}

func (buf *Spill) Cap() int64 {
	return MAXINT64
}

func (buf *Spill) Write(p []byte) (n int, err error) {
	if n, err = buf.Buffer.Write(p); err != nil {
		buf.Spiller.Write(p[:n])
	}
	return len(p), nil
}

func (buf *Spill) WriteAt(p []byte, off int64) (n int, err error) {
	if n, err = buf.Buffer.WriteAt(p, off); err != nil {
		buf.Spiller.Write(p[:n])
	}
	return len(p), nil
}

func init() {
	gob.Register(&Spill{})
}
