package buffer

type Queue interface {
	Peek() []byte
	Pop() []byte
	Push([]byte)
	Len() int64
}

type SliceQueue [][]byte

func NewSliceQueue() Queue {
	q := SliceQueue(make([][]byte, 0))
	return &q
}

func (q *SliceQueue) Push(b []byte) {
	*q = append(*q, b)
}

func (q *SliceQueue) Peek() (b []byte) {
	if len(*q) > 0 {
		b = (*q)[0]
	}
	return b
}

func (q *SliceQueue) Pop() (b []byte) {
	if len(*q) > 0 {
		b = (*q)[0]
		*q = (*q)[1:]
	}
	return b
}

func (q *SliceQueue) Len() int64 {
	return int64(len(*q))
}

type BufferQueue struct {
	top []byte
	Buffer
}

func NewBufferQueue(buf Buffer) Queue {
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

func (q *BufferQueue) Peek() (b []byte) {
	if q.top != nil {
		return q.top
	} else {
		q.top = q.Pop()
		return q.top
	}
}

func (q *BufferQueue) Push(b []byte) {
	q.Write(b)
}

func (q *BufferQueue) Pop() (b []byte) {
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
