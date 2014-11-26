package buffer

import (
	"encoding/gob"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

var dataDir string
var setDataDir sync.Once

// BUG(Dustin): Path is created even if its not used.
func SetDataDir(path string, perm os.FileMode) (err error) {
	if err = os.MkdirAll(path, perm); err == nil {
		err = errors.New("already set")
		setDataDir.Do(func() {
			dataDir = path
			err = nil
		})
	}
	return err
}

type File struct {
	file     *os.File
	Filename string
	Wrapper
}

func NewFile(N int64) *File {
	buf := &File{}
	buf.Wrapper.N = N
	buf.init()
	return buf
}

func (buf *File) init() (err error) {
	var file *os.File

	if buf.file == nil {

		setDataDir.Do(func() {
			dataDir = ""
		})

		if buf.Filename != "" {
			if file, err = os.OpenFile(buf.Filename, os.O_CREATE|os.O_RDWR, 0644); err != nil {
				return err
			}
		} else if file, err = ioutil.TempFile(dataDir, "buffer"); err != nil {
			return err
		}

		buf.file = file
		buf.Filename = file.Name()
		buf.Wrapper.rwa = file
	}

	return nil
}

func (buf *File) Read(p []byte) (n int, err error) {
	if err = buf.init(); err != nil {
		return n, err
	}
	return buf.Wrapper.Read(p)
}

func (buf *File) ReadAt(p []byte, off int64) (n int, err error) {
	if err = buf.init(); err != nil {
		return n, err
	}
	return buf.Wrapper.ReadAt(p, off)
}

func (buf *File) Write(p []byte) (n int, err error) {
	if err = buf.init(); err != nil {
		return n, err
	}
	return buf.Wrapper.Write(p)
}

func (buf *File) WriteAt(p []byte, off int64) (n int, err error) {
	if err = buf.init(); err != nil {
		return n, err
	}
	return buf.Wrapper.WriteAt(p, off)
}

func (buf *File) Reset() {
	if buf.file != nil {
		buf.file.Close()
		os.Remove(buf.file.Name())
		buf.file = nil
		buf.Filename = ""
		buf.Wrapper.L = 0
		buf.Wrapper.O = 0
	}
}

func init() {
	gob.Register(&File{})
}
