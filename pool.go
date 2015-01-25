package buffer

import (
	"encoding/gob"
	"io/ioutil"
	"os"
	"sync"
)

// Pool provides a way to Allocate and Release Buffer objects
type Pool interface {
	Get() Buffer    // Allocate a Buffer
	Put(buf Buffer) // Release or Reuse a Buffer
}

type pool struct {
	pool sync.Pool
}

// NewPool returns a Pool(), it's backed by a sync.Pool so its safe for concurrent use.
// Unlike NewMemPool or NewFilePool, this pool supports concurrent calls to Get() and Put().
// However it will not work with gob.
func NewPool(New func() Buffer) Pool {
	return &pool{
		pool: sync.Pool{
			New: func() interface{} {
				return New()
			},
		},
	}
}

func (p *pool) Get() Buffer {
	return p.pool.Get().(Buffer)
}

func (p *pool) Put(buf Buffer) {
	buf.Reset()
	p.pool.Put(buf)
}

type memPool struct {
	N int64
}

// NewMemPool returns a Pool, Get() returns an in memory buffer of max size N.
// Put() is a nop.
func NewMemPool(N int64) Pool {
	return &memPool{N: N}
}

func (p *memPool) Get() Buffer    { return New(p.N) }
func (p *memPool) Put(buf Buffer) {}

type filePool struct {
	N         int64
	Directory string
}

// NewFilePool returns a Pool, Get() returns a file-based buffer of max size N.
// Put() closes and deletes the underlying file for the buffer.
func NewFilePool(N int64, dir string) Pool {
	return &filePool{N: N, Directory: dir}
}

func (p *filePool) Get() Buffer {
	file, err := ioutil.TempFile(p.Directory, "buffer")
	if err != nil {
		panic(err.Error())
	}
	return NewFile(p.N, file)
}

func (p *filePool) Put(buf Buffer) {
	if fileBuf, ok := buf.(*fileBuffer); ok {
		fileBuf.file.Close()
		os.Remove(fileBuf.file.Name())
	}
}

func init() {
	gob.Register(&memPool{})
	gob.Register(&filePool{})
}
