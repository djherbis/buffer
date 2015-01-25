package buffer

import (
	"io"
	"os"
)

const MAXINT64 = 9223372036854775807

type Buffer interface {
	Len() int64
	Cap() int64
	io.Reader
	io.Writer
	Reset()
}

type BufferAt interface {
	Buffer
	io.ReaderAt
	io.WriterAt
}

func len64(p []byte) int64 {
	return int64(len(p))
}

// Gap returns buf.Cap() - buf.Len()
func Gap(buf Buffer) int64 {
	return buf.Cap() - buf.Len()
}

// Full returns true iff buf.Len() == buf.Cap()
func Full(buf Buffer) bool {
	return buf.Len() == buf.Cap()
}

// Empty returns false iff buf.Len() == 0
func Empty(buf Buffer) bool {
	return buf.Len() == 0
}

// NewUnboundedBuffer returns a Buffer which buffers "mem" bytes to memory
// and then creates file's of size "file" to buffer above "mem" bytes.
func NewUnboundedBuffer(mem, file int64) Buffer {
	return NewMulti(New(mem), NewPartition(NewFilePool(file, os.TempDir())))
}
