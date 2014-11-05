package buffer

import (
	"io"
	"io/ioutil"
	"os"
)

type FileBuffer struct {
	file        *os.File
	cap         int64
	len         int64
	readOffset  int64
	writeOffset int64
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
			buf.readOffset = 0
			buf.writeOffset = 0
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

func (buf *FileBuffer) ReadAt(p []byte, off int64) (n int, err error) {
	wrap := NewWrapReader(buf.file, off, buf.Cap())
	r := io.NewSectionReader(wrap, 0, buf.Len()-off)
	return r.ReadAt(p, 0)
}

func (buf *FileBuffer) Read(p []byte) (n int, err error) {
	wrap := NewWrapReader(buf.file, buf.readOffset, buf.Cap())
	r := io.LimitReader(wrap, buf.Len())

	n, err = r.Read(p)
	buf.len -= int64(n)
	buf.readOffset = (buf.readOffset + int64(n)) % buf.Cap()

	if Empty(buf) {
		buf.Reset()
	}

	return n, err
}

// BUG(Dustin): Add short write
func (buf *FileBuffer) Write(p []byte) (n int, err error) {
	if err := buf.init(); err != nil {
		return 0, err
	}

	wrap := NewWrapWriter(buf.file, buf.writeOffset, buf.Cap())
	w := LimitWriter(wrap, Gap(buf))
	n, err = w.Write(p)
	buf.len += int64(n)
	buf.writeOffset = (buf.writeOffset + int64(n)) % buf.Cap()

	return n, err
}

func (buf *FileBuffer) Reset() {
	if buf.file != nil {
		buf.file.Close()
		os.Remove(buf.file.Name())
		buf.file = nil
		buf.readOffset = 0
		buf.writeOffset = 0
	}
}

func (buf *FileBuffer) FastForward(n int) int {
	if int64(n) > buf.Len() {
		n = int(buf.Len())
	}

	buf.len -= int64(n)
	buf.readOffset = (buf.readOffset + int64(n)) % buf.Cap()

	return n
}
