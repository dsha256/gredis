package cache

import (
	"time"
)

// DataType represents the type of data stored in the cache.
type DataType int

const (
	// StringType represents a string value.
	StringType DataType = iota
	// ListType represents a list value.
	ListType
)

// StringCmdable defines the interface for string operations.
type StringCmdable interface {
	Get(key string) (string, bool)
	Set(key string, value string) error
	SetWithTTL(key string, value string, ttl time.Duration) error
	Update(key string, value string) error
}

// ListCmdable defines the interface for list operations.
type ListCmdable interface {
	PushFront(key string, value string) error
	PushBack(key string, value string) error
	PopFront(key string) (string, bool)
	PopBack(key string) (string, bool)
	ListRange(key string, start, end int) ([]string, error)
}

// TTLCmdable defines the interface for TTL operations.
type TTLCmdable interface {
	SetTTL(key string, ttl time.Duration) error
	GetTTL(key string) (time.Duration, bool)
	RemoveTTL(key string) error
}

// GeneralCmdable defines the interface for general operations.
type GeneralCmdable interface {
	Remove(key string) error
	Exists(key string) bool
	Type(key string) (DataType, bool)
	Clear() error
}

// Cache defines the interface for all cache operations.
type Cache interface {
	StringCmdable
	ListCmdable
	TTLCmdable
	GeneralCmdable
}
