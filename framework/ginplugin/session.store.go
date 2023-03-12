package ginplugin

import (
	"sync"
	"time"
)

type SessionData map[string][]byte

type SessionStore interface {
	Load(id string) (SessionData, error)
	Save(id string, sessionData SessionData, ttlInMillis int) error
	TouchIfExists(id string, ttlInMillis int) error
}

type memSessionStore struct {
	sync.RWMutex

	data map[string]entry
}

type entry struct {
	val SessionData
	ttl int64
}

func (m *memSessionStore) Load(id string) (SessionData, error) {
	m.RLock()
	val, ok := m.data[id]
	m.RUnlock()
	now := time.Now().UnixMilli()
	if !ok {
		return nil, nil
	} else if val.ttl <= now {
		return nil, nil
	} else {
		return val.val, nil
	}
}

func (m *memSessionStore) Save(id string, sessionData SessionData, ttlInMillis int) error {
	ttl := time.Now().UnixMilli() + int64(ttlInMillis)
	m.Lock()
	val, ok := m.data[id]
	if !ok {
		val = entry{}
	}
	val.ttl = ttl
	val.val = sessionData
	m.data[id] = val
	m.Unlock()
	return nil
}

func (m *memSessionStore) TouchIfExists(id string, ttlInMillis int) error {
	ttl := time.Now().UnixMilli() + int64(ttlInMillis)
	m.Lock()
	val, ok := m.data[id]
	if ok {
		val.ttl = ttl
		m.data[id] = val
	}
	m.Unlock()
	return nil
}

func (m *memSessionStore) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		<-ticker.C
		now := time.Now().UnixMilli()
		func() {
			m.Lock()
			defer m.Unlock()

			for k, v := range m.data {
				if v.ttl <= now {
					delete(m.data, k)
				}
			}
		}()
	}
}

func NewMemSessionStore() SessionStore {
	store := &memSessionStore{
		data: map[string]entry{},
	}
	go store.cleanupLoop()
	return store
}
