package buffer

type BufferQueue struct {
	top []byte
	Buffer
}

func NewBufferQueue(buf Buffer) *BufferQueue {
	return &BufferQueue{
		Buffer: buf,
	}
}

func (q *BufferQueue) Len() (n int64) {
	if q.top != nil {
		n += int64(len(q.top))
	}
	n += q.Buffer.Len()
	return n
}

func (q *BufferQueue) Peek() (b interface{}) {
	if q.top != nil {
		return q.top
	} else {
		q.top = q.Pop().([]byte)
		return q.top
	}
}

func (q *BufferQueue) Push(b interface{}) {
	q.Write(b.([]byte))
}

func (q *BufferQueue) Pop() (b interface{}) {
	if q.top != nil {
		b = q.top
		q.top = nil
	} else {
		data := make([]byte, 32*1024)
		n, _ := q.Read(data)
		b = data[:n]
	}
	return b
}
