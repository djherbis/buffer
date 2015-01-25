package buffer

import (
	"io"

	"github.com/djherbis/buffer/wrapio"
)

type File interface {
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
}

type FileBuffer struct {
	file File
	N    int64
	*wrapio.Wrapper
}

func NewFile(N int64, file File) *FileBuffer {
	return &FileBuffer{
		file:    file,
		N:       N,
		Wrapper: wrapio.NewWrapper(file, N),
	}
}

func (buf *FileBuffer) Reset() {
	buf.Wrapper.L = 0
	buf.Wrapper.O = 0
}
