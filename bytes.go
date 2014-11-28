package buffer

type Bytes []byte

func NewBytes(n int64) BufferAt {
	return NewWrapper(Bytes(make([]byte, n)), n)
}

// BUG(Dustin): off past end?
func (b Bytes) WriteAt(p []byte, off int64) (n int, err error) {
	return copy(b[off:], p), nil
}

// BUG(Dustin): off past end?
func (b Bytes) ReadAt(p []byte, off int64) (n int, err error) {
	return copy(p, b[off:]), nil
}

// BUG(Dustin): doesn't support gob yet.
