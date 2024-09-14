package qexmap

import (
	"sync"
	"time"
)

/*
过期map
特点：
1.设置key时必须设置key的存在时间
2.加锁的安全map

优化：过期key的清理
1.在获取key时，判断距离上次清理时间超过60秒则启动协程清理一次
2.无需启动一个永久协程循环检查清理失效的key，浪费资源
*/

var cleanTime time.Time

type ExpiringMap struct {
	mu    sync.RWMutex
	store map[string]item
}

type item struct {
	value      interface{}
	expiration time.Time
}

// NewExpiringMap 创建一个新的过期 map
func NewExpiringMap() *ExpiringMap {
	return &ExpiringMap{
		store: make(map[string]item),
	}
}

// Set 设置键值以及过期时间
func (m *ExpiringMap) Set(key string, value interface{}, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] = item{
		value:      value,
		expiration: time.Now().Add(duration),
	}
}

// Get 获取键的值，如果过期则返回 nil
func (m *ExpiringMap) Get(key string) (interface{}, bool) {
	if time.Now().After(cleanTime) {
		go m.CleanUp()
		cleanTime = time.Now().Add(60 * time.Second)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	itm, found := m.store[key]
	if !found || time.Now().After(itm.expiration) {
		return nil, false
	}
	return itm.value, true
}

// CleanUp 定期清理过期的键
func (m *ExpiringMap) CleanUp() {
	go func() {
		m.mu.Lock()
		for key, itm := range m.store {
			if time.Now().After(itm.expiration) {
				delete(m.store, key)
			}
		}
		m.mu.Unlock()
	}()
}
