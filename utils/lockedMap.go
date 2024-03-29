package utils

import "sync"

type LockedMap struct {
	m map[string]interface{}
	sync.RWMutex
}

func NewLockedMap() *LockedMap {
	return &LockedMap{
		m:       make(map[string]interface{}),
		RWMutex: sync.RWMutex{},
	}
}

func (m *LockedMap) Get(key string) (interface{}, bool) {
	m.RLock()
	defer m.RUnlock()
	v, ok := m.m[key]
	return v, ok
}

func (m *LockedMap) Put(key string, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.m[key] = value
}

func (m *LockedMap) Remove(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, key)
}

func (m *LockedMap) Iter(handler func(key string, value interface{})) {
	m.Lock()
	defer m.Unlock()
	for key, value := range m.m {
		handler(key, value)
	}
}

func (m *LockedMap) Size() int {
	m.Lock()
	defer m.Unlock()
	return len(m.m)
}
