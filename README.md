# Gredis

Gredis is a simple in-memory data structure store inspired by Redis, implemented in Go.

## Features

- **Data Structures**:
  - Strings
  - Lists

- **Operations**:
  - Get
  - Set
  - Update
  - Remove
  - Push for lists (PushFront, PushBack)
  - Pop for lists (PopFront, PopBack)

- **Additional Features**:
  - Keys with a limited TTL (Time To Live)
  - Go client API library
  - Automatic cleanup of expired keys

## Installation

```bash
go get github.com/dsha256/gredis
```

## Usage

### Basic Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/dsha256/gredis/client"
)

func main() {
	// Create a new memory client with a cleanup interval of 1 second
	c := client.NewMemoryClient(1 * time.Second)
	defer c.Close() // Make sure to close the client when done

	// Set a string value
	c.Set("greeting", "Hello, World!")

	// Get the value
	value, err := c.Get("greeting")
	if err == nil {
		fmt.Printf("greeting = %s\n", value)
	}

	// Set a value with TTL
	c.SetWithTTL("temp", "This will expire", 5*time.Second)

	// Create a list
	c.PushBack("mylist", "first")
	c.PushBack("mylist", "second")
	c.PushFront("mylist", "zero")

	// Get the list range
	items, err := c.ListRange("mylist", 0, -1)
	if err == nil {
		fmt.Printf("List items: %v\n", items)
	}
}
```

### Using Specialized Clients

Gredis provides specialized clients for different types of operations:

- `StringClient` for string operations
- `ListClient` for list operations
- `TTLClient` for TTL operations

```go
package main

import (
	"fmt"
	"time"

	"github.com/dsha256/gredis/client"
)

func main() {
	// Create a new memory client with a cleanup interval of 1 second
	c := client.NewMemoryClient(1 * time.Second)
	defer c.Close() // Make sure to close the client when done

	// Get specialized clients
	strClient := c.String()
	listClient := c.List()
	ttlClient := c.TTL()

	// Use StringClient for string operations
	strClient.Set("greeting", "Hello, World!")
	value, _ := strClient.Get("greeting")

	// Use ListClient for list operations
	listClient.PushBack("mylist", "first")
	listClient.PushFront("mylist", "zero")
	items, _ := listClient.ListRange("mylist", 0, -1)

	// Use TTLClient for TTL operations
	strClient.SetWithTTL("temp", "This will expire", 5*time.Second)
	ttl, _ := ttlClient.GetTTL("temp")
}
```

### String Operations

Using the main client:

```go
// Set a string value
c.Set("key", "value")

// Set a string value with TTL
c.SetWithTTL("key", "value", 10*time.Second)

// Get a string value
value, err := c.Get("key")

// Update an existing string value
c.Update("key", "new value")

// Remove a key
c.Remove("key")
```

Using the specialized StringClient:

```go
// Get the StringClient
strClient := c.String()

// Set a string value
strClient.Set("key", "value")

// Set a string value with TTL
strClient.SetWithTTL("key", "value", 10*time.Second)

// Get a string value
value, err := strClient.Get("key")

// Update an existing string value
strClient.Update("key", "new value")
```

### List Operations

Using the main client:

```go
// Add to the front of a list
c.PushFront("list", "value")

// Add to the back of a list
c.PushBack("list", "value")

// Remove and return the first element of a list
value, err := c.PopFront("list")

// Remove and return the last element of a list
value, err := c.PopBack("list")

// Get a range of elements from a list
// Use negative indices to count from the end (-1 is the last element)
items, err := c.ListRange("list", 0, -1) // Get all elements
```

Using the specialized ListClient:

```go
// Get the ListClient
listClient := c.List()

// Add to the front of a list
listClient.PushFront("list", "value")

// Add to the back of a list
listClient.PushBack("list", "value")

// Remove and return the first element of a list
value, err := listClient.PopFront("list")

// Remove and return the last element of a list
value, err := listClient.PopBack("list")

// Get a range of elements from a list
items, err := listClient.ListRange("list", 0, -1) // Get all elements
```

### TTL Operations

Using the main client:

```go
// Set TTL for a key
c.SetWithTTL("key", "value", 10*time.Second)

// Get the remaining TTL for a key
ttl, err := c.GetTTL("key")

// Remove TTL for a key (make it persistent)
c.RemoveTTL("key")
```

Using the specialized TTLClient:

```go
// Get the TTLClient
ttlClient := c.TTL()

// Set TTL for a key (requires StringClient)
strClient := c.String()
strClient.SetWithTTL("key", "value", 10*time.Second)

// Get the remaining TTL for a key
ttl, err := ttlClient.GetTTL("key")

// Remove TTL for a key (make it persistent)
ttlClient.RemoveTTL("key")
```

### Other Operations

```go
// Check if a key exists
exists := c.Exists("key")

// Get the type of a key
dataType, err := c.Type("key")

// Clear all keys
c.Clear()

// Close the client when done
c.Close()
```

## Examples

See the [examples](examples) directory for more detailed examples.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
