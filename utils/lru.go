package utils

import (
	"container/list"
	"sync"
)

// LruMap is a thread-safe map with a LRU (Least Recently Used) eviction policy.
// When the cache reaches its maximum size, the least recently used item is evicted.
//
// See: https://en.wikipedia.org/wiki/Cache_replacement_policies#LRU
type LruMap[V any] struct {
	maxSize int
	cache   map[string]*list.Element
	order   *list.List
	mu      sync.RWMutex
}

// lruEntry stores the key-value pair in the LRU list.
type lruEntry[V any] struct {
	key   string
	value V
}

// NewLruMap creates a new LRU map with the specified maximum size.
//
// Example:
//
//	cache := NewLruMap[int](100)
//	cache.Set("key1", 42)
//	value, ok := cache.Get("key1")
func NewLruMap[V any](maxSize int) *LruMap[V] {
	return &LruMap[V]{
		maxSize: maxSize,
		cache:   make(map[string]*list.Element),
		order:   list.New(),
	}
}

// Get retrieves a value from the cache.
// If the key exists, it's moved to the front (most recently used).
// Returns the value and true if found, zero value and false otherwise.
func (l *LruMap[V]) Get(key string) (V, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if elem, ok := l.cache[key]; ok {
		l.order.MoveToFront(elem)
		return elem.Value.(*lruEntry[V]).value, true
	}

	var zero V
	return zero, false
}

// Set adds or updates a value in the cache.
// If the key already exists, its value is updated and moved to the front.
// If adding a new key would exceed the max size, the least recently used item is evicted.
func (l *LruMap[V]) Set(key string, value V) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if elem, ok := l.cache[key]; ok {
		// Update existing entry
		elem.Value.(*lruEntry[V]).value = value
		l.order.MoveToFront(elem)
		return
	}

	// Add new entry
	entry := &lruEntry[V]{key: key, value: value}
	elem := l.order.PushFront(entry)
	l.cache[key] = elem

	// Evict oldest if over capacity
	if l.maxSize > 0 && l.order.Len() > l.maxSize {
		l.evictOldest()
	}
}

// Delete removes a key from the cache.
// Returns true if the key was present and removed.
func (l *LruMap[V]) Delete(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if elem, ok := l.cache[key]; ok {
		l.removeElement(elem)
		return true
	}
	return false
}

// Has checks if a key exists in the cache without affecting its LRU position.
func (l *LruMap[V]) Has(key string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	_, ok := l.cache[key]
	return ok
}

// Size returns the current number of items in the cache.
func (l *LruMap[V]) Size() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.order.Len()
}

// MaxSize returns the maximum size of the cache.
func (l *LruMap[V]) MaxSize() int {
	return l.maxSize
}

// Clear removes all items from the cache.
func (l *LruMap[V]) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cache = make(map[string]*list.Element)
	l.order.Init()
}

// Keys returns all keys in the cache, ordered from most to least recently used.
func (l *LruMap[V]) Keys() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	keys := make([]string, 0, l.order.Len())
	for elem := l.order.Front(); elem != nil; elem = elem.Next() {
		keys = append(keys, elem.Value.(*lruEntry[V]).key)
	}
	return keys
}

// evictOldest removes the least recently used item (must be called with lock held).
func (l *LruMap[V]) evictOldest() {
	oldest := l.order.Back()
	if oldest != nil {
		l.removeElement(oldest)
	}
}

// removeElement removes an element from the cache (must be called with lock held).
func (l *LruMap[V]) removeElement(elem *list.Element) {
	l.order.Remove(elem)
	entry := elem.Value.(*lruEntry[V])
	delete(l.cache, entry.key)
}
