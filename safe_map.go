package safemap

import (
	"sync"
	"sync/atomic"
)

//Автор - Чурин Владимир Михаилович

type SafeMap struct {
	mu        sync.Mutex
	data      map[int]*int64
	accesses  int
	additions int
}

func New() *SafeMap {
	return &SafeMap{
		data: make(map[int]*int64),
	}
}

func (m *SafeMap) GetOrCreate(key int) *int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.accesses++

	if v, ok := m.data[key]; ok {
		return v
	}

	v := new(int64)
	m.data[key] = v
	m.additions++
	return v
}

func (m *SafeMap) Load(key int) (int64, bool) {
	m.mu.Lock()
	v, ok := m.data[key]
	m.mu.Unlock()
	if !ok {
		return 0, false
	}
	return atomic.LoadInt64(v), true
}

func (m *SafeMap) Stats() (accesses int, additions int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.accesses, m.additions
}
