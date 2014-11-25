package buffer

import (
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
