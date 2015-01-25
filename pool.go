package buffer

import (
	"encoding/gob"
	"io/ioutil"
	"os"
	"sync"
)

type Pool struct {
	pool sync.Pool
}

func NewPool(New func() Buffer) *Pool {
	return &Pool{
		pool: sync.Pool{
			New: func() interface{} {
				return New()
			},
		},
	}
}

func (p *Pool) Get() Buffer {
	return p.pool.Get().(Buffer)
}

func (p *Pool) Put(buf Buffer) {
	buf.Reset()
	p.pool.Put(buf)
}

type memPool struct {
	N int64
}

func NewMemPool(N int64) *memPool {
	return &memPool{N: N}
}

func (p *memPool) Get() Buffer    { return New(p.N) }
func (p *memPool) Put(buf Buffer) {}

type filePool struct {
	N int64
}

func NewFilePool(N int64) *filePool {
	return &filePool{N: N}
}

func (p *filePool) Get() Buffer {
	file, err := ioutil.TempFile("", "buffer")
	if err != nil {
		panic(err.Error())
	}
	return NewFile(p.N, file)
}

func (p *filePool) Put(buf Buffer) {
	if fileBuf, ok := buf.(*FileBuffer); ok {
		fileBuf.file.Close()
		os.Remove(fileBuf.file.Name())
	}
}

func init() {
	gob.Register(&memPool{})
	gob.Register(&filePool{})
}
