package buffer

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"os"
	"sync"

	"github.com/djherbis/buffer/wrapio"
)

var dataDir string
var setDataDir sync.Once

func DataDir() string {
	if dataDir == "" {
		return os.TempDir()
	}
	return dataDir
}

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
	N        int64
	*wrapio.Wrapper
}

func NewFile(N int64) *File {
	buf := &File{
		N:       N,
		Wrapper: wrapio.NewWrapper(nil, N),
	}
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
		buf.Wrapper.SetReadWriterAt(file)
	}

	return nil
}

func (buf *File) Read(p []byte) (n int, err error) {
	if err = buf.init(); err != nil {
		return n, err
	}
	n, err = buf.Wrapper.Read(p)
	if Empty(buf) {
		buf.Reset()
	}
	return n, err
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

func (buf *File) MarshalBinary() ([]byte, error) {
	w := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(w)
	enc.Encode(&buf.Filename)
	enc.Encode(&buf.N)
	enc.Encode(&buf.Wrapper)
	return w.Bytes(), nil
}

func (buf *File) UnmarshalBinary(data []byte) (err error) {
	r := bytes.NewBuffer(data)
	dec := gob.NewDecoder(r)
	dec.Decode(&buf.Filename)
	dec.Decode(&buf.N)
	dec.Decode(&buf.Wrapper)
	buf.init()
	return err
}

func init() {
	gob.Register(&File{})
}
