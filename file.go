package buffer

import (
	"encoding/gob"
	"io"
	"io/ioutil"
	"os"
)

type File struct {
	file     *os.File
	Filename string
	N        int64
	L        int64
	ROff     int64
	WOff     int64
}

func NewFile(N int64) *File {
	buf := &File{
		N: N,
	}
	return buf
}

func (buf *File) init() (err error) {
	var file *os.File

	if buf.file == nil {

		if buf.Filename != "" {
			if file, err = os.OpenFile(buf.Filename, os.O_CREATE|os.O_RDWR, 0644); err != nil {
				return err
			}
		} else if file, err = ioutil.TempFile("", "buffer"); err != nil {
			return err
		}

		buf.file = file
		buf.Filename = file.Name()
		buf.ROff = 0
		buf.WOff = 0
	}

	return nil
}

func (buf *File) Len() int64 {
	return buf.L
}

func (buf *File) Cap() int64 {
	return buf.N
}

func (buf *File) ReadAt(p []byte, off int64) (n int, err error) {
	if err := buf.init(); err != nil {
		return 0, err
	}

	wrap := NewWrapReader(buf.file, buf.ROff+off, buf.Cap())
	r := io.LimitReader(wrap, buf.Len()-off)
	return r.Read(p)
}

func (buf *File) Read(p []byte) (n int, err error) {
	if err := buf.init(); err != nil {
		return 0, err
	}

	wrap := NewWrapReader(buf.file, buf.ROff, buf.Cap())
	r := io.LimitReader(wrap, buf.Len())

	n, err = r.Read(p)
	buf.L -= int64(n)
	buf.ROff = (buf.ROff + int64(n)) % buf.Cap()

	if Empty(buf) {
		buf.Reset()
	}

	return n, err
}

func (buf *File) Write(p []byte) (n int, err error) {
	if err := buf.init(); err != nil {
		return 0, err
	}

	wrap := NewWrapWriter(buf.file, buf.WOff, buf.Cap())
	w := LimitWriter(wrap, Gap(buf))
	n, err = w.Write(p)
	buf.L += int64(n)
	buf.WOff = (buf.WOff + int64(n)) % buf.Cap()

	return n, err
}

func (buf *File) WriteAt(p []byte, off int64) (n int, err error) {
	if err := buf.init(); err != nil {
		return 0, err
	}

	wrap := NewWrapWriter(buf.file, off, buf.Cap())
	w := LimitWriter(wrap, buf.Cap()-off)
	n, err = w.Write(p)
	if off+int64(n) > buf.Len() {
		buf.L = int64(n) + off
	}
	return n, err
}

func (buf *File) Reset() {
	if buf.file != nil {
		buf.file.Close()
		os.Remove(buf.file.Name())
		buf.file = nil
		buf.Filename = ""
		buf.ROff = 0
		buf.WOff = 0
	}
}

func (buf *File) FFwd(n int64) int64 {
	if int64(n) > buf.Len() {
		n = buf.Len()
	}

	buf.L -= int64(n)
	buf.ROff = (buf.ROff + int64(n)) % buf.Cap()

	return n
}

func init() {
	gob.Register(&File{})
}
