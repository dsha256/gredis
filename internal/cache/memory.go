package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// Common errors
var (
	ErrKeyNotFound  = errors.New("key not found")
	ErrTypeMismatch = errors.New("type mismatch")
)

// cacheItem represents a value stored in the cache
type cacheItem struct {
	dataType DataType
	value    any
	expireAt time.Time // Zero time means no expiration
}

// isExpired checks if the item has expired
func (i *cacheItem) isExpired() bool {
	return !i.expireAt.IsZero() && time.Now().After(i.expireAt)
}

// MemoryCache implements the Cache interface with in-memory storage
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
	// For TTL cleanup
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache(cleanupInterval time.Duration) *MemoryCache {
	cache := &MemoryCache{
		items:           make(map[string]*cacheItem),
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan struct{}),
	}

	// Start cleanup goroutine if interval is positive
	if cleanupInterval > 0 {
		go cache.startCleanup()
	}

	return cache
}

// startCleanup starts the cleanup process for expired items
func (c *MemoryCache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

// cleanup removes expired items
func (c *MemoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if !item.expireAt.IsZero() && now.After(item.expireAt) {
			delete(c.items, key)
		}
	}
}

// Stop stops the cleanup goroutine
func (c *MemoryCache) Stop() {
	if c.cleanupInterval > 0 {
		c.stopCleanup <- struct{}{}
	}
}

// Get retrieves a string value from the cache
func (c *MemoryCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || item.isExpired() {
		if found && item.isExpired() {
			// Cleanup expired item..
			c.mu.RUnlock()
			c.mu.Lock()
			delete(c.items, key)
			c.mu.Unlock()
			c.mu.RLock()
		}
		return "", false
	}

	if item.dataType != StringType {
		return "", false
	}

	return item.value.(string), true
}

// Set stores a string value in the cache
func (c *MemoryCache) Set(key string, value string) error {
	return c.set(key, value, 0)
}

// SetWithTTL stores a string value in the cache with a TTL
func (c *MemoryCache) SetWithTTL(key string, value string, ttl time.Duration) error {
	return c.set(key, value, ttl)
}

// set is a helper function for Set and SetWithTTL
func (c *MemoryCache) set(key string, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expireAt time.Time
	if ttl > 0 {
		expireAt = time.Now().Add(ttl)
	}

	c.items[key] = &cacheItem{
		dataType: StringType,
		value:    value,
		expireAt: expireAt,
	}

	return nil
}

// Update updates an existing string value in the cache
func (c *MemoryCache) Update(key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found || item.isExpired() {
		if found && item.isExpired() {
			delete(c.items, key)
		}
		return ErrKeyNotFound
	}

	if item.dataType != StringType {
		return ErrTypeMismatch
	}

	item.value = value
	return nil
}

// Remove removes a key from the cache
func (c *MemoryCache) Remove(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.items[key]
	if !found {
		return ErrKeyNotFound
	}

	delete(c.items, key)
	return nil
}

// PushFront adds a value to the front of a list.
func (c *MemoryCache) PushFront(key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		// Create a new list. if the key doesn't exist..
		l := list.New()
		l.PushFront(value)
		c.items[key] = &cacheItem{
			dataType: ListType,
			value:    l,
			expireAt: time.Time{},
		}
		return nil
	}

	if item.isExpired() {
		delete(c.items, key)
		// Create a new list..
		l := list.New()
		l.PushFront(value)
		c.items[key] = &cacheItem{
			dataType: ListType,
			value:    l,
			expireAt: time.Time{},
		}
		return nil
	}

	if item.dataType != ListType {
		return ErrTypeMismatch
	}

	l := item.value.(*list.List)
	l.PushFront(value)
	return nil
}

// PushBack adds a value to the back of a list.
func (c *MemoryCache) PushBack(key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		// Create a new list. if the key doesn't exist..
		l := list.New()
		l.PushBack(value)
		c.items[key] = &cacheItem{
			dataType: ListType,
			value:    l,
			expireAt: time.Time{},
		}
		return nil
	}

	if item.isExpired() {
		delete(c.items, key)
		// Create a new list..
		l := list.New()
		l.PushBack(value)
		c.items[key] = &cacheItem{
			dataType: ListType,
			value:    l,
			expireAt: time.Time{},
		}
		return nil
	}

	if item.dataType != ListType {
		return ErrTypeMismatch
	}

	l := item.value.(*list.List)
	l.PushBack(value)
	return nil
}

// PopFront removes and returns the first element of a list.
func (c *MemoryCache) PopFront(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found || item.isExpired() {
		if found && item.isExpired() {
			delete(c.items, key)
		}
		return "", false
	}

	if item.dataType != ListType {
		return "", false
	}

	l := item.value.(*list.List)
	if l.Len() == 0 {
		return "", false
	}

	element := l.Front()
	l.Remove(element)
	return element.Value.(string), true
}

// PopBack removes and returns the last element of a list.
func (c *MemoryCache) PopBack(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found || item.isExpired() {
		if found && item.isExpired() {
			delete(c.items, key)
		}
		return "", false
	}

	if item.dataType != ListType {
		return "", false
	}

	l := item.value.(*list.List)
	if l.Len() == 0 {
		return "", false
	}

	element := l.Back()
	l.Remove(element)
	return element.Value.(string), true
}

// ListRange returns a range of elements from a list.
func (c *MemoryCache) ListRange(key string, start, end int) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || item.isExpired() {
		if found && item.isExpired() {
			// Cleanup expired item..
			c.mu.RUnlock()
			c.mu.Lock()
			delete(c.items, key)
			c.mu.Unlock()
			c.mu.RLock()
		}
		return nil, ErrKeyNotFound
	}

	if item.dataType != ListType {
		return nil, ErrTypeMismatch
	}

	l := item.value.(*list.List)
	length := l.Len()

	// Handle negative indices..
	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}

	// Validate indices..
	if start < 0 {
		start = 0
	}
	if end >= length {
		end = length - 1
	}
	if start > end || start >= length {
		return []string{}, nil
	}

	// Extract the range..
	result := make([]string, 0, end-start+1)
	e := l.Front()
	for i := 0; i < start; i++ {
		e = e.Next()
	}
	for i := start; i <= end; i++ {
		result = append(result, e.Value.(string))
		e = e.Next()
	}

	return result, nil
}

// SetTTL sets the TTL for a key.
func (c *MemoryCache) SetTTL(key string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found || item.isExpired() {
		if found && item.isExpired() {
			delete(c.items, key)
		}
		return ErrKeyNotFound
	}

	if ttl <= 0 {
		item.expireAt = time.Time{}
	} else {
		item.expireAt = time.Now().Add(ttl)
	}

	return nil
}

// GetTTL returns the remaining TTL for a key.
func (c *MemoryCache) GetTTL(key string) (time.Duration, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || item.isExpired() {
		if found && item.isExpired() {
			// Cleanup expired item..
			c.mu.RUnlock()
			c.mu.Lock()
			delete(c.items, key)
			c.mu.Unlock()
			c.mu.RLock()
		}
		return 0, false
	}

	if item.expireAt.IsZero() {
		return -1, true // -1 indicates no expiration...
	}

	ttl := time.Until(item.expireAt)
	if ttl < 0 {
		return 0, false
	}

	return ttl, true
}

// RemoveTTL removes the TTL for a key.
func (c *MemoryCache) RemoveTTL(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found || item.isExpired() {
		if found && item.isExpired() {
			delete(c.items, key)
		}
		return ErrKeyNotFound
	}

	item.expireAt = time.Time{}
	return nil
}

// Exists checks if a key exists in the cache.
func (c *MemoryCache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return false
	}

	if item.isExpired() {
		// Cleanup expired item..
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		c.mu.RLock()
		return false
	}

	return true
}

// Type returns the type of a key.
func (c *MemoryCache) Type(key string) (DataType, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || item.isExpired() {
		if found && item.isExpired() {
			// Cleanup expired item..
			c.mu.RUnlock()
			c.mu.Lock()
			delete(c.items, key)
			c.mu.Unlock()
			c.mu.RLock()
		}
		return 0, false
	}

	return item.dataType, true
}

// Clear removes all items from the cache.
func (c *MemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
	return nil
}
