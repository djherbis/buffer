package buffer

import (
	"encoding/gob"
	"io"
	"io/ioutil"
	"os"
)

type FileBuffer struct {
	file     *os.File
	Filename string
	N        int64
	L        int64
	ROff     int64
	WOff     int64
}

func NewFile(N int64) *FileBuffer {
	buf := &FileBuffer{
		N: N,
	}
	return buf
}

func (buf *FileBuffer) init() (err error) {
	var file *os.File

	if buf.file == nil {

		if buf.Filename != "" {
			if file, err = os.OpenFile(buf.Filename, os.O_CREATE|os.O_RDWR, 0644); err != nil {
				return err
			}
		} else if file, err = ioutil.TempFile("D:\\Downloads\\temp", "buffer"); err != nil {
			return err
		}

		buf.file = file
		buf.Filename = file.Name()
		buf.ROff = 0
		buf.WOff = 0
	}

	return nil
}

func (buf *FileBuffer) Len() int64 {
	return buf.L
}

func (buf *FileBuffer) Cap() int64 {
	return buf.N
}

func (buf *FileBuffer) ReadAt(p []byte, off int64) (n int, err error) {
	if err := buf.init(); err != nil {
		return 0, err
	}

	wrap := NewWrapReader(buf.file, off, buf.Cap())
	r := io.NewSectionReader(wrap, 0, buf.Len()-off)
	return r.ReadAt(p, 0)
}

func (buf *FileBuffer) Read(p []byte) (n int, err error) {
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

func (buf *FileBuffer) Write(p []byte) (n int, err error) {
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

func (buf *FileBuffer) Reset() {
	if buf.file != nil {
		buf.file.Close()
		os.Remove(buf.file.Name())
		buf.file = nil
		buf.Filename = ""
		buf.ROff = 0
		buf.WOff = 0
	}
}

func (buf *FileBuffer) FFwd(n int) int {
	if int64(n) > buf.Len() {
		n = int(buf.Len())
	}

	buf.L -= int64(n)
	buf.ROff = (buf.ROff + int64(n)) % buf.Cap()

	return n
}

func init() {
	gob.Register(&FileBuffer{})
}
