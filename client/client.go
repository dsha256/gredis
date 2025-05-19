package client

import (
	"errors"
	"time"

	"github.com/dsha256/gredis/internal/cache"
)

// Common errors.
var (
	ErrKeyNotFound        = errors.New("key not found")
	ErrKeyNotFoundOrEmpty = errors.New("key not found or empty list")
)

// Client provides a client API for interacting with the cache.
type Client struct {
	cache cache.Cache
}

// StringClient provides a client API for string operations.
type StringClient struct {
	cmdable cache.StringCmdable
}

// ListClient provides a client API for list operations.
type ListClient struct {
	cmdable cache.ListCmdable
}

// New creates a new client with the given cache implementation.
func New(cache cache.Cache) *Client {
	return &Client{
		cache: cache,
	}
}

// NewMemoryClient creates a new client with an in-memory cache.
func NewMemoryClient(cleanupInterval time.Duration) *Client {
	return &Client{
		cache: cache.NewMemoryCache(cleanupInterval),
	}
}

// String returns a client for string operations.
func (c *Client) String() *StringClient {
	return &StringClient{
		cmdable: c.cache,
	}
}

// List returns a client for list operations.
func (c *Client) List() *ListClient {
	return &ListClient{
		cmdable: c.cache,
	}
}

// String operations.

// Get retrieves a string value from the cache.
func (c *StringClient) Get(key string) (string, error) {
	value, ok := c.cmdable.Get(key)
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}

// Set stores a string value in the cache.
func (c *StringClient) Set(key string, value string) error {
	return c.cmdable.Set(key, value)
}

// SetWithTTL stores a string value in the cache with a TTL.
func (c *StringClient) SetWithTTL(key string, value string, ttl time.Duration) error {
	return c.cmdable.SetWithTTL(key, value, ttl)
}

// Update updates an existing string value in the cache.
func (c *StringClient) Update(key string, value string) error {
	return c.cmdable.Update(key, value)
}

// Remove removes a key from the cache.
func (c *Client) Remove(key string) error {
	return c.cache.Remove(key)
}

// List operations.

// PushFront adds a value to the front of a list.
func (c *ListClient) PushFront(key string, value string) error {
	return c.cmdable.PushFront(key, value)
}

// PushBack adds a value to the back of a list.
func (c *ListClient) PushBack(key string, value string) error {
	return c.cmdable.PushBack(key, value)
}

// PopFront removes and returns the first element of a list.
func (c *ListClient) PopFront(key string) (string, error) {
	value, ok := c.cmdable.PopFront(key)
	if !ok {
		return "", ErrKeyNotFoundOrEmpty
	}
	return value, nil
}

// PopBack removes and returns the last element of a list.
func (c *ListClient) PopBack(key string) (string, error) {
	value, ok := c.cmdable.PopBack(key)
	if !ok {
		return "", ErrKeyNotFoundOrEmpty
	}
	return value, nil
}

// ListRange returns a range of elements from a list.
func (c *ListClient) ListRange(key string, start, end int) ([]string, error) {
	return c.cmdable.ListRange(key, start, end)
}

// Get retrieves a string value from the cache.
func (c *Client) Get(key string) (string, error) {
	return c.String().Get(key)
}

// Set stores a string value in the cache.
func (c *Client) Set(key string, value string) error {
	return c.String().Set(key, value)
}

// SetWithTTL stores a string value in the cache with a TTL.
func (c *Client) SetWithTTL(key string, value string, ttl time.Duration) error {
	return c.String().SetWithTTL(key, value, ttl)
}

// Update updates an existing string value in the cache.
func (c *Client) Update(key string, value string) error {
	return c.String().Update(key, value)
}

// PushFront adds a value to the front of a list.
func (c *Client) PushFront(key string, value string) error {
	return c.List().PushFront(key, value)
}

// PushBack adds a value to the back of a list.
func (c *Client) PushBack(key string, value string) error {
	return c.List().PushBack(key, value)
}

// PopFront removes and returns the first element of a list.
func (c *Client) PopFront(key string) (string, error) {
	return c.List().PopFront(key)
}

// PopBack removes and returns the last element of a list.
func (c *Client) PopBack(key string) (string, error) {
	return c.List().PopBack(key)
}

// ListRange returns a range of elements from a list.
func (c *Client) ListRange(key string, start, end int) ([]string, error) {
	return c.List().ListRange(key, start, end)
}

// TTLClient provides a client API for TTL operations.
type TTLClient struct {
	cmdable cache.TTLCmdable
}

// TTL returns a client for TTL operations.
func (c *Client) TTL() *TTLClient {
	return &TTLClient{
		cmdable: c.cache,
	}
}

// TTL operations.

// SetTTL sets the TTL for a key.
func (c *TTLClient) SetTTL(key string, ttl time.Duration) error {
	return c.cmdable.SetTTL(key, ttl)
}

// GetTTL returns the remaining TTL for a key.
func (c *TTLClient) GetTTL(key string) (time.Duration, error) {
	ttl, ok := c.cmdable.GetTTL(key)
	if !ok {
		return 0, ErrKeyNotFound
	}
	return ttl, nil
}

// RemoveTTL removes the TTL for a key.
func (c *TTLClient) RemoveTTL(key string) error {
	return c.cmdable.RemoveTTL(key)
}

// SetTTL sets the TTL for a key.
func (c *Client) SetTTL(key string, ttl time.Duration) error {
	return c.TTL().SetTTL(key, ttl)
}

// GetTTL returns the remaining TTL for a key.
func (c *Client) GetTTL(key string) (time.Duration, error) {
	return c.TTL().GetTTL(key)
}

// RemoveTTL removes the TTL for a key.
func (c *Client) RemoveTTL(key string) error {
	return c.TTL().RemoveTTL(key)
}

// General operations.

// Exists checks if a key exists in the cache.
func (c *Client) Exists(key string) bool {
	return c.cache.Exists(key)
}

// Type returns the type of a key.
func (c *Client) Type(key string) (cache.DataType, error) {
	dataType, ok := c.cache.Type(key)
	if !ok {
		return 0, ErrKeyNotFound
	}
	return dataType, nil
}

// Clear removes all items from the cache.
func (c *Client) Clear() error {
	return c.cache.Clear()
}

// Close closes the client and releases any resources.
func (c *Client) Close() error {
	if memCache, ok := c.cache.(*cache.MemoryCache); ok {
		memCache.Stop()
	}
	return nil
}
