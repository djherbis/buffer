package buffer

import "sync"

type strset struct {
	sync.RWMutex
	objects map[string]struct{}
}

var tracker strset = strset{
	objects: make(map[string]struct{}),
}

func (t *strset) Track(s string) {
	t.Lock()
	defer t.Unlock()
	t.objects[s] = struct{}{}
}

func (t *strset) Untrack(s string) {
	t.Lock()
	defer t.Unlock()
	delete(t.objects, s)
}

func (t *strset) IsTracked(s string) bool {
	t.RLock()
	defer t.RUnlock()
	_, ok := t.objects[s]
	return ok
}
