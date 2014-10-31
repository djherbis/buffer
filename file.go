package buffer

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

type file struct {
	file *os.File
	cap  int64
	len  int64
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
}

func NewFile(cap int64) Buffer {
	buf := &file{
		cap: cap,
	}
	return buf
}

func (buf *file) init() error {
	if buf.file == nil {
		if file, err := ioutil.TempFile("D:\\Downloads\\temp", "buffer"); err == nil {
			buf.file = file
			r := NewWrapReader(file, buf.Cap())
			buf.Reader = r
			buf.ReaderAt = r
			w := NewWrapWriter(file, buf.Cap())
			buf.Writer = w
			buf.WriterAt = w
		} else {
			return err
		}
	}
	return nil
}

func (buf *file) Len() int64 {
	return buf.len
}

func (buf *file) Cap() int64 {
	return buf.cap
}

func (buf *file) Read(p []byte) (n int, err error) {
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
func (buf *file) Write(p []byte) (n int, err error) {
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

func (buf *file) Reset() {
	if buf.file != nil {
		buf.file.Close()
		os.Remove(buf.file.Name())
		buf.file = nil
		buf.Reader = nil
		buf.Writer = nil
	}
}
