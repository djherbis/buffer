package buffer

import (
	"io"
	"sync"
)

type Sync struct {
	done chan struct{}
	l    sync.Mutex
	c    *sync.Cond
	b    Buffer
	lr   int
}

func NewSync(buf Buffer) *Sync {
	s := &Sync{
		b:    buf,
		done: make(chan struct{}),
	}
	s.c = sync.NewCond(&s.l)
	return s
}

func (r *Sync) Done() {
	close(r.done)
	r.c.Signal()
}

func (r *Sync) Read(p []byte) (n int, err error) {
	r.l.Lock()
	defer r.c.Signal()
	defer r.l.Unlock()

	r.b.FFwd(int64(r.lr))

	for Empty(r.b) {
		select {
		case <-r.done:
			return 0, io.EOF
		default:
		}

		r.c.Signal()
		r.c.Wait()
	}

	n, _ = r.b.ReadAt(p, 0)
	r.lr = n

	return n, nil
}

func (w *Sync) Write(p []byte) (n int, err error) {
	w.l.Lock()
	defer w.c.Signal()
	defer w.l.Unlock()

	for Gap(w.b) < len64(p) {
		w.c.Signal()
		w.c.Wait()
	}

	return w.b.Write(p)
}
