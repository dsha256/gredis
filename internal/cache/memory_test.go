package cache

import (
	"errors"
	"testing"
	"time"
)

func TestMemoryCache_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(c *MemoryCache)
		key        string
		value      string
		operation  string
		wantValue  string
		wantExists bool
		wantErr    error
	}{
		{
			name:       "Get non-existent key",
			setup:      func(c *MemoryCache) {},
			key:        "nonexistent",
			operation:  "Get",
			wantValue:  "",
			wantExists: false,
		},
		{
			name: "Get existing key",
			setup: func(c *MemoryCache) {
				err := c.Set("key1", "value1")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:        "key1",
			operation:  "Get",
			wantValue:  "value1",
			wantExists: true,
		},
		{
			name:       "Set new key",
			setup:      func(c *MemoryCache) {},
			key:        "key2",
			value:      "value2",
			operation:  "Set",
			wantExists: true,
		},
		{
			name: "Update existing key",
			setup: func(c *MemoryCache) {
				err := c.Set("key3", "value3")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "key3",
			value:     "updated3",
			operation: "Update",
			wantErr:   nil,
		},
		{
			name:      "Update non-existent key",
			setup:     func(c *MemoryCache) {},
			key:       "nonexistent",
			value:     "value",
			operation: "Update",
			wantErr:   ErrKeyNotFound,
		},
		{
			name: "Get expired key",
			setup: func(c *MemoryCache) {
				err := c.SetWithTTL("expired", "value", 1*time.Millisecond)
				requireNoError(t, err, "Setup failed: %v", err)
				time.Sleep(10 * time.Millisecond) // Ensure key expires
			},
			key:        "expired",
			operation:  "Get",
			wantValue:  "",
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewMemoryCache(100 * time.Millisecond)
			defer c.Stop()

			tt.setup(c)

			switch tt.operation {
			case "Get":
				gotValue, gotExists := c.Get(tt.key)
				require(t, gotExists == tt.wantExists, "Get() exists = %v, want %v", gotExists, tt.wantExists)
				if gotExists {
					require(t, gotValue == tt.wantValue, "Get() value = %v, want %v", gotValue, tt.wantValue)
				}
			case "Set":
				err := c.Set(tt.key, tt.value)
				require(t, errors.Is(err, nil), "Set() error = %v", err)

				// Verify the set worked
				gotValue, gotExists := c.Get(tt.key)
				require(t, gotExists, "Set() key not found after setting")
				require(t, gotValue == tt.value, "Set() value = %v, want %v", gotValue, tt.value)
			case "Update":
				err := c.Update(tt.key, tt.value)
				require(t, errors.Is(err, tt.wantErr), "Update() error = %v, want %v", err, tt.wantErr)

				if errors.Is(err, nil) {
					// Verify the update worked
					gotValue, gotExists := c.Get(tt.key)
					require(t, gotExists, "Update() key not found after update")
					require(t, gotValue == tt.value, "Update() value = %v, want %v", gotValue, tt.value)
				}
			}
		})
	}
}

func TestMemoryCache_TTL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(c *MemoryCache)
		key        string
		ttl        time.Duration
		operation  string
		wantTTL    time.Duration
		wantExists bool
		wantErr    error
	}{
		{
			name: "SetTTL on existing key",
			setup: func(c *MemoryCache) {
				err := c.Set("key1", "value1")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "key1",
			ttl:       5 * time.Second,
			operation: "SetTTL",
			wantErr:   nil,
		},
		{
			name:      "SetTTL on non-existent key",
			setup:     func(c *MemoryCache) {},
			key:       "nonexistent",
			ttl:       5 * time.Second,
			operation: "SetTTL",
			wantErr:   ErrKeyNotFound,
		},
		{
			name: "GetTTL on key with TTL",
			setup: func(c *MemoryCache) {
				err := c.Set("key2", "value2")
				requireNoError(t, err, "Setup failed: %v", err)
				err = c.SetTTL("key2", 5*time.Second)
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:        "key2",
			operation:  "GetTTL",
			wantExists: true,
		},
		{
			name: "GetTTL on key without TTL",
			setup: func(c *MemoryCache) {
				err := c.Set("key3", "value3")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:        "key3",
			operation:  "GetTTL",
			wantTTL:    -1, // -1 indicates no expiration
			wantExists: true,
		},
		{
			name:       "GetTTL on non-existent key",
			setup:      func(c *MemoryCache) {},
			key:        "nonexistent",
			operation:  "GetTTL",
			wantExists: false,
		},
		{
			name: "RemoveTTL on key with TTL",
			setup: func(c *MemoryCache) {
				err := c.Set("key4", "value4")
				requireNoError(t, err, "Setup failed: %v", err)
				err = c.SetTTL("key4", 5*time.Second)
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "key4",
			operation: "RemoveTTL",
			wantErr:   nil,
		},
		{
			name:      "RemoveTTL on non-existent key",
			setup:     func(c *MemoryCache) {},
			key:       "nonexistent",
			operation: "RemoveTTL",
			wantErr:   ErrKeyNotFound,
		},
		{
			name: "GetTTL on expired key",
			setup: func(c *MemoryCache) {
				err := c.Set("expired", "value")
				requireNoError(t, err, "Setup failed: %v", err)
				err = c.SetTTL("expired", 1*time.Millisecond)
				requireNoError(t, err, "Setup failed: %v", err)
				time.Sleep(10 * time.Millisecond) // Ensure key expires
			},
			key:        "expired",
			operation:  "GetTTL",
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewMemoryCache(100 * time.Millisecond)
			defer c.Stop()

			tt.setup(c)

			switch tt.operation {
			case "SetTTL":
				err := c.SetTTL(tt.key, tt.ttl)
				require(t, errors.Is(err, tt.wantErr), "SetTTL() error = %v, want %v", err, tt.wantErr)

				if errors.Is(err, nil) {
					// Verify the TTL was set
					ttl, exists := c.GetTTL(tt.key)
					require(t, exists, "SetTTL() key not found after setting TTL")
					// We can't check exact TTL as it depends on execution time, but we can check it's positive
					require(t, ttl > 0 || ttl == -1, "SetTTL() TTL not set correctly")
				}

			case "GetTTL":
				ttl, exists := c.GetTTL(tt.key)
				require(t, exists == tt.wantExists, "GetTTL() exists = %v, want %v", exists, tt.wantExists)

				if exists && tt.wantTTL == -1 {
					require(t, ttl == -1, "GetTTL() TTL = %v, want %v", ttl, tt.wantTTL)
				}

			case "RemoveTTL":
				err := c.RemoveTTL(tt.key)
				require(t, errors.Is(err, tt.wantErr), "RemoveTTL() error = %v, want %v", err, tt.wantErr)

				if errors.Is(err, nil) {
					// Verify the TTL was removed
					ttl, exists := c.GetTTL(tt.key)
					require(t, exists, "RemoveTTL() key not found after removing TTL")
					require(t, ttl == -1, "RemoveTTL() TTL not removed correctly")
				}
			}
		})
	}
}

func TestMemoryCache_General(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(c *MemoryCache)
		key        string
		operation  string
		wantExists bool
		wantType   DataType
		wantErr    error
	}{
		{
			name: "Remove existing key",
			setup: func(c *MemoryCache) {
				err := c.Set("key1", "value1")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "key1",
			operation: "Remove",
			wantErr:   nil,
		},
		{
			name:      "Remove non-existent key",
			setup:     func(c *MemoryCache) {},
			key:       "nonexistent",
			operation: "Remove",
			wantErr:   ErrKeyNotFound,
		},
		{
			name: "Exists with existing key",
			setup: func(c *MemoryCache) {
				err := c.Set("key2", "value2")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:        "key2",
			operation:  "Exists",
			wantExists: true,
		},
		{
			name:       "Exists with non-existent key",
			setup:      func(c *MemoryCache) {},
			key:        "nonexistent",
			operation:  "Exists",
			wantExists: false,
		},
		{
			name: "Type with string key",
			setup: func(c *MemoryCache) {
				err := c.Set("string", "value")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "string",
			operation: "Type",
			wantType:  StringType,
		},
		{
			name: "Type with list key",
			setup: func(c *MemoryCache) {
				err := c.PushBack("list", "value")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "list",
			operation: "Type",
			wantType:  ListType,
		},
		{
			name:      "Type with non-existent key",
			setup:     func(c *MemoryCache) {},
			key:       "nonexistent",
			operation: "Type",
			wantErr:   ErrKeyNotFound,
		},
		{
			name: "Clear cache",
			setup: func(c *MemoryCache) {
				err := c.Set("key1", "value1")
				requireNoError(t, err, "Setup failed: %v", err)
				err = c.Set("key2", "value2")
				requireNoError(t, err, "Setup failed: %v", err)
				err = c.PushBack("list", "value")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			operation: "Clear",
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewMemoryCache(100 * time.Millisecond)
			defer c.Stop()

			tt.setup(c)

			switch tt.operation {
			case "Remove":
				err := c.Remove(tt.key)
				require(t, errors.Is(err, tt.wantErr), "Remove() error = %v, want %v", err, tt.wantErr)

				if errors.Is(err, nil) {
					// Verify the key was removed
					exists := c.Exists(tt.key)
					require(t, !exists, "Remove() key still exists after removal")
				}

			case "Exists":
				exists := c.Exists(tt.key)
				require(t, exists == tt.wantExists, "Exists() = %v, want %v", exists, tt.wantExists)

			case "Type":
				dataType, exists := c.Type(tt.key)

				if errors.Is(tt.wantErr, ErrKeyNotFound) {
					require(t, !exists, "Type() exists = %v, want false", exists)
				} else {
					require(t, exists, "Type() exists = %v, want true", exists)
					require(t, dataType == tt.wantType, "Type() = %v, want %v", dataType, tt.wantType)
				}

			case "Clear":
				err := c.Clear()
				require(t, errors.Is(err, tt.wantErr), "Clear() error = %v, want %v", err, tt.wantErr)

				// Verify all keys were removed
				exists := c.Exists("key1") || c.Exists("key2") || c.Exists("list")
				require(t, !exists, "Clear() keys still exist after clearing")
			}
		})
	}
}

func TestMemoryCache_List(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(c *MemoryCache)
		key       string
		value     string
		operation string
		start     int
		end       int
		wantValue string
		wantList  []string
		wantOk    bool
		wantErr   error
	}{
		{
			name:      "PushFront to new list",
			setup:     func(c *MemoryCache) {},
			key:       "list1",
			value:     "value1",
			operation: "PushFront",
			wantErr:   nil,
		},
		{
			name: "PushFront to existing list",
			setup: func(c *MemoryCache) {
				err := c.PushBack("list2", "value2")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "list2",
			value:     "value1",
			operation: "PushFront",
			wantErr:   nil,
		},
		{
			name:      "PushBack to new list",
			setup:     func(c *MemoryCache) {},
			key:       "list3",
			value:     "value3",
			operation: "PushBack",
			wantErr:   nil,
		},
		{
			name: "PushBack to existing list",
			setup: func(c *MemoryCache) {
				err := c.PushFront("list4", "value4")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "list4",
			value:     "value5",
			operation: "PushBack",
			wantErr:   nil,
		},
		{
			name:      "PopFront from empty list",
			setup:     func(c *MemoryCache) {},
			key:       "emptylist",
			operation: "PopFront",
			wantValue: "",
			wantOk:    false,
		},
		{
			name: "PopFront from list with one item",
			setup: func(c *MemoryCache) {
				err := c.PushBack("list5", "value5")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "list5",
			operation: "PopFront",
			wantValue: "value5",
			wantOk:    true,
		},
		{
			name:      "PopBack from empty list",
			setup:     func(c *MemoryCache) {},
			key:       "emptylist",
			operation: "PopBack",
			wantValue: "",
			wantOk:    false,
		},
		{
			name: "PopBack from list with one item",
			setup: func(c *MemoryCache) {
				err := c.PushBack("list6", "value6")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "list6",
			operation: "PopBack",
			wantValue: "value6",
			wantOk:    true,
		},
		{
			name: "ListRange full list",
			setup: func(c *MemoryCache) {
				err := c.PushBack("list7", "value1")
				requireNoError(t, err, "Setup failed: %v", err)
				err = c.PushBack("list7", "value2")
				requireNoError(t, err, "Setup failed: %v", err)
				err = c.PushBack("list7", "value3")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "list7",
			operation: "ListRange",
			start:     0,
			end:       -1,
			wantList:  []string{"value1", "value2", "value3"},
			wantErr:   nil,
		},
		{
			name: "ListRange partial list",
			setup: func(c *MemoryCache) {
				err := c.PushBack("list8", "value1")
				requireNoError(t, err, "Setup failed: %v", err)
				err = c.PushBack("list8", "value2")
				requireNoError(t, err, "Setup failed: %v", err)
				err = c.PushBack("list8", "value3")
				requireNoError(t, err, "Setup failed: %v", err)
				err = c.PushBack("list8", "value4")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "list8",
			operation: "ListRange",
			start:     1,
			end:       2,
			wantList:  []string{"value2", "value3"},
			wantErr:   nil,
		},
		{
			name:      "ListRange non-existent key",
			setup:     func(c *MemoryCache) {},
			key:       "nonexistent",
			operation: "ListRange",
			start:     0,
			end:       -1,
			wantList:  nil,
			wantErr:   ErrKeyNotFound,
		},
		{
			name: "ListRange type mismatch",
			setup: func(c *MemoryCache) {
				err := c.Set("string1", "value1")
				requireNoError(t, err, "Setup failed: %v", err)
			},
			key:       "string1",
			operation: "ListRange",
			start:     0,
			end:       -1,
			wantList:  nil,
			wantErr:   ErrTypeMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewMemoryCache(100 * time.Millisecond)
			defer c.Stop()

			tt.setup(c)

			switch tt.operation {
			case "PushFront":
				err := c.PushFront(tt.key, tt.value)
				require(t, errors.Is(err, tt.wantErr), "PushFront() error = %v, want %v", err, tt.wantErr)

				// Verify the push worked
				items, err := c.ListRange(tt.key, 0, 0)
				require(t, errors.Is(err, nil), "PushFront() verification failed: %v", err)
				require(t, len(items) > 0 && items[0] == tt.value, "PushFront() value not found in list")

			case "PushBack":
				err := c.PushBack(tt.key, tt.value)
				require(t, errors.Is(err, tt.wantErr), "PushBack() error = %v, want %v", err, tt.wantErr)

				// Verify the push worked
				items, err := c.ListRange(tt.key, -1, -1)
				require(t, errors.Is(err, nil), "PushBack() verification failed: %v", err)
				require(t, len(items) > 0 && items[0] == tt.value, "PushBack() value not found in list")

			case "PopFront":
				gotValue, gotOk := c.PopFront(tt.key)
				require(t, gotOk == tt.wantOk, "PopFront() ok = %v, want %v", gotOk, tt.wantOk)
				if gotOk {
					require(t, gotValue == tt.wantValue, "PopFront() value = %v, want %v", gotValue, tt.wantValue)
				}

			case "PopBack":
				gotValue, gotOk := c.PopBack(tt.key)
				require(t, gotOk == tt.wantOk, "PopBack() ok = %v, want %v", gotOk, tt.wantOk)
				if gotOk {
					require(t, gotValue == tt.wantValue, "PopBack() value = %v, want %v", gotValue, tt.wantValue)
				}

			case "ListRange":
				gotList, err := c.ListRange(tt.key, tt.start, tt.end)
				require(t, errors.Is(err, tt.wantErr), "ListRange() error = %v, want %v", err, tt.wantErr)

				if errors.Is(err, nil) {
					require(t, len(gotList) == len(tt.wantList), "ListRange() list length = %v, want %v", len(gotList), len(tt.wantList))

					for i, v := range gotList {
						require(t, v == tt.wantList[i], "ListRange() list[%d] = %v, want %v", i, v, tt.wantList[i])
					}
				}
			}
		})
	}
}

func requireNoError(t *testing.T, err error, format string, args ...any) {
	t.Helper()
	require(t, errors.Is(err, nil), format, args...)
}

func require(t *testing.T, condition bool, format string, args ...any) {
	t.Helper()
	if !condition {
		t.Fatalf(format, args...)
	}
}
