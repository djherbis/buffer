package buffer

import (
	"io"
	"io/ioutil"
	"os"
)

type file struct {
	file        *os.File
	capacity    int64
	readOffset  int64
	writeOffset int64
}

func NewFile(capacity int64) Buffer {
	return &file{
		capacity: capacity,
	}
}

func (buf *file) init() error {
	if buf.file == nil {
		if file, err := ioutil.TempFile("", "buffer"); err == nil {
			buf.file = file
		} else {
			return err
		}
	}
	return nil
}

func (buf *file) Len() int64 {
	return buf.writeOffset - buf.readOffset
}

func (buf *file) Cap() int64 {
	return buf.capacity
}

func (buf *file) Read(p []byte) (n int, err error) {
	if Empty(buf) {
		return 0, io.EOF
	}

	if err := buf.init(); err != nil {
		return 0, err
	}

	buf.file.Seek(buf.readOffset, 0)
	n, err = buf.file.Read(p)
	buf.readOffset += int64(n)

	if Empty(buf) {
		buf.remove()
	}

	return n, err
}

func (buf *file) Write(p []byte) (n int, err error) {
	if Full(buf) {
		return 0, err
	}

	if err := buf.init(); err != nil {
		return 0, err
	}

	buf.file.Seek(buf.writeOffset, 0)
	n, err = buf.file.Write(ShrinkToFit(buf, p))
	buf.writeOffset += int64(n)

	return n, err
}

func (buf *file) remove() {
	buf.file.Close()
	os.Remove(buf.file.Name())
	buf.file = nil
}
