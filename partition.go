package buffer

import (
	"encoding/gob"
	"errors"
	"io"
	"math"
)

type partition struct {
	List
	Pool
}

// NewPartition returns a Buffer which uses a Pool to extend or shrink its size as needed.
// It automatically allocates new buffers with pool.Get() to extend is length, and
// pool.Put() to release unused buffers as it shrinks.
func NewPartition(pool Pool, buffers ...Buffer) Buffer {
	return &partition{
		Pool: pool,
		List: buffers,
	}
}

func (buf *partition) Cap() int64 {
	return math.MaxInt64
}

func (buf *partition) Read(p []byte) (n int, err error) {
	for len(p) > 0 {

		if len(buf.List) == 0 {
			return n, io.EOF
		}

		buffer := buf.List[0]

		if Empty(buffer) {
			buf.Pool.Put(buf.Pop())
			continue
		}

		m, er := buffer.Read(p)
		n += m
		p = p[m:]

		if er != nil && er != io.EOF {
			return n, er
		}

	}
	return n, nil
}

func (buf *partition) Write(p []byte) (n int, err error) {
	for len(p) > 0 {

		if len(buf.List) == 0 {
			next, err := buf.Pool.Get()
			if err != nil {
				return n, err
			}
			buf.Push(next)
		}

		buffer := buf.List[len(buf.List)-1]

		if Full(buffer) {
			next, err := buf.Pool.Get()
			if err != nil {
				return n, err
			}
			buf.Push(next)
			continue
		}

		m, er := buffer.Write(p)
		n += m
		p = p[m:]

		if er == io.ErrShortWrite {
			er = nil
		} else if er != nil {
			return n, er
		}

	}
	return n, nil
}

func (buf *partition) Reset() {
	for len(buf.List) > 0 {
		buf.Pool.Put(buf.Pop())
	}
}

func init() {
	gob.Register(&partition{})
}

type partitionAt struct {
	ListAt
	PoolAt
}

// NewPartitionAt returns a BufferAt which uses a PoolAt to extend or shrink its size as needed.
// It automatically allocates new buffers with pool.Get() to extend is length, and
// pool.Put() to release unused buffers as it shrinks.
func NewPartitionAt(pool PoolAt, buffers ...BufferAt) BufferAt {
	return &partitionAt{
		PoolAt: pool,
		ListAt: buffers,
	}
}

func (buf *partitionAt) Cap() int64 {
	return math.MaxInt64
}

func (buf *partitionAt) Read(p []byte) (n int, err error) {
	for len(p) > 0 {

		if len(buf.ListAt) == 0 {
			return n, io.EOF
		}

		buffer := buf.ListAt[0]

		if Empty(buffer) {
			buf.PoolAt.Put(buf.Pop())
			continue
		}

		m, er := buffer.Read(p)
		n += m
		p = p[m:]

		if er != nil && er != io.EOF {
			return n, er
		}

	}
	return n, nil
}

func (buf *partitionAt) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("buffer.PartionAt.ReadAt: negative offset")
	}
	for _, buffer := range buf.ListAt {

		m, er := buffer.ReadAt(p, off)
		n += m
		p = p[m:]

		if er != nil && er != io.EOF {
			return n, er
		}
		if len(p) == 0 {
			return n, er
		}
		off = off - buffer.Len()
		if off < 0 {
			off = 0
		}
	}
	if len(p) > 0 {
		err = io.EOF
	}
	return n, err
}

func (buf *partitionAt) Write(p []byte) (n int, err error) {
	for len(p) > 0 {

		if len(buf.ListAt) == 0 {
			next, err := buf.PoolAt.Get()
			if err != nil {
				return n, err
			}
			buf.Push(next)
		}

		buffer := buf.ListAt[len(buf.ListAt)-1]

		if Full(buffer) {
			next, err := buf.PoolAt.Get()
			if err != nil {
				return n, err
			}
			buf.Push(next)
			continue
		}

		m, er := buffer.Write(p)
		n += m
		p = p[m:]

		if er == io.ErrShortWrite {
			er = nil
		} else if er != nil {
			return n, er
		}

	}
	return n, nil
}

func (buf *partitionAt) WriteAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("buffer.PartionAt.WriteAt: negative offset")
	}
	for _, buffer := range buf.ListAt {

		m, er := buffer.WriteAt(p, off)
		n += m
		p = p[m:]

		if er != nil && er != io.ErrShortWrite {
			return n, er
		}
		if len(p) == 0 {
			return n, er
		}
		off = off - buffer.Len()
		if off < 0 {
			off = 0
		}
	}
	if len(p) > 0 {
		err = io.ErrShortWrite
	}
	return n, err
}

func (buf *partitionAt) Reset() {
	for len(buf.ListAt) > 0 {
		buf.PoolAt.Put(buf.Pop())
	}
}

func init() {
	gob.Register(&partitionAt{})
}
