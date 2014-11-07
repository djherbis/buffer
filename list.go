package buffer

type BufferList []Buffer

func (l *BufferList) Len() (n int64) {
	for _, buffer := range *l {
		if n > MAXINT64-buffer.Len() {
			return MAXINT64
		} else {
			n += buffer.Len()
		}
	}
	return n
}

func (l *BufferList) Cap() (n int64) {
	for _, buffer := range *l {
		if n > MAXINT64-buffer.Cap() {
			return MAXINT64
		} else {
			n += buffer.Cap()
		}
	}
	return n
}

func (l *BufferList) Reset() {
	for _, buffer := range *l {
		buffer.Reset()
	}
}

func (l *BufferList) Push(b Buffer) {
	*l = append(*l, b)
}

func (l *BufferList) Pop() (b Buffer) {
	b = (*l)[0]
	*l = (*l)[1:]
	return b
}
