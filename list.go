package buffer

type BufferList []Buffer

func (l *BufferList) Len() int64 {
	return TotalLen(*l)
}

func (l *BufferList) Cap() int64 {
	return TotalCap(*l)
}

func (l *BufferList) Reset() {
	ResetAll(*l)
}

func (l *BufferList) Push(b Buffer) {
	*l = append(*l, b)
}

func (l *BufferList) Pop() (b Buffer) {
	b = (*l)[0]
	*l = (*l)[1:]
	return b
}
