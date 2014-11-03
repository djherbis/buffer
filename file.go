package buffer

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

type FileBuffer struct {
	file *os.File
	cap  int64
	len  int64
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
}

func NewFile(cap int64) *FileBuffer {
	buf := &FileBuffer{
		cap: cap,
	}
	return buf
}

func (buf *FileBuffer) init() error {
	if buf.file == nil {
		if file, err := ioutil.TempFile("D:\\Downloads\\temp", "buffer"); err == nil {
			buf.file = file
			r := NewWrapReader(file, 0, buf.Cap())
			buf.Reader = r
			buf.ReaderAt = r
			w := NewWrapWriter(file, 0, buf.Cap())
			buf.Writer = w
			buf.WriterAt = w
		} else {
			return err
		}
	}
	return nil
}

func (buf *FileBuffer) Len() int64 {
	return buf.len
}

func (buf *FileBuffer) Cap() int64 {
	return buf.cap
}

func (buf *FileBuffer) Read(p []byte) (n int, err error) {
	if Empty(buf) {
		return 0, io.EOF
	}

	n, err = buf.Reader.Read(ShrinkToRead(buf, p))
	buf.len -= int64(n)

	if Empty(buf) {
		buf.Reset()
	}

	return n, err
}

// BUG(Dustin): Add short write
func (buf *FileBuffer) Write(p []byte) (n int, err error) {
	if Full(buf) {
		return 0, bytes.ErrTooLarge
	}

	if err := buf.init(); err != nil {
		return 0, err
	}

	n, err = buf.Writer.Write(ShrinkToFit(buf, p))
	buf.len += int64(n)

	return n, err
}

func (buf *FileBuffer) Reset() {
	if buf.file != nil {
		buf.file.Close()
		os.Remove(buf.file.Name())
		buf.file = nil
		buf.Reader = nil
		buf.Writer = nil
	}
}
