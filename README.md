# Gredis 🚀

Gredis is a simple in-memory data structure store inspired by Redis, implemented in Go.

## Table of Contents 📑
- [Features](#features-)
- [Installation](#installation)
- [Usage](#usage)
  - [Basic Usage](#basic-usage)
  - [Using Specialized Clients](#using-specialized-clients)
  - [String Operations](#string-operations)
  - [List Operations](#list-operations)
  - [TTL Operations](#ttl-operations)
  - [Other Operations](#other-operations)
- [API Endpoints](#api-endpoints-)
  - [String Operations](#string-operations-api)
  - [List Operations](#list-operations-api)
  - [TTL Operations](#ttl-operations-api)
  - [General Operations](#general-operations-api)
- [Running Locally with Docker](#running-locally-with-docker-)
  - [Using Docker Directly](#using-docker-directly)
  - [Using Docker Compose](#using-docker-compose)
  - [Using Taskfile](#using-taskfile)
- [Taskfile Commands](#taskfile-commands)
- [Examples](#examples)
- [License](#license)

## Features ✨

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

## API Endpoints 🌐

Gredis provides a RESTful API for interacting with the cache. Below are the available endpoints and examples of how to use them with cURL.

### String Operations API

#### Get a string value

```
GET /api/v1/string/{key}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8090/api/v1/string/greeting
```

**Response:**
```json
{
  "data": {
    "key": "greeting",
    "value": "Hello, World!"
  },
  "msg": "Value retrieved successfully"
}
```

#### Set a string value

```
POST /api/v1/string/{key}
```

**Request Body:**
```json
{
  "value": "Hello, World!",
  "ttl": 60
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8090/api/v1/string/greeting \
  -H "Content-Type: application/json" \
  -d '{"value": "Hello, World!", "ttl": 60}'
```

**Response:**
```json
{
  "data": {
    "key": "greeting",
    "value": "Hello, World!"
  },
  "msg": "Value set successfully"
}
```

#### Update a string value

```
PUT /api/v1/string/{key}
```

**Request Body:**
```json
{
  "value": "Updated value"
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8090/api/v1/string/greeting \
  -H "Content-Type: application/json" \
  -d '{"value": "Updated value"}'
```

**Response:**
```json
{
  "data": {
    "key": "greeting",
    "value": "Updated value"
  },
  "msg": "Value updated successfully"
}
```

### List Operations API

#### Push a value to the front of a list

```
POST /api/v1/list/{key}/front
```

**Request Body:**
```json
{
  "value": "first item"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8090/api/v1/list/mylist/front \
  -H "Content-Type: application/json" \
  -d '{"value": "first item"}'
```

**Response:**
```json
{
  "data": {
    "key": "mylist",
    "value": "first item"
  },
  "msg": "Value pushed to front of list successfully"
}
```

#### Push a value to the back of a list

```
POST /api/v1/list/{key}/back
```

**Request Body:**
```json
{
  "value": "last item"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8090/api/v1/list/mylist/back \
  -H "Content-Type: application/json" \
  -d '{"value": "last item"}'
```

**Response:**
```json
{
  "data": {
    "key": "mylist",
    "value": "last item"
  },
  "msg": "Value pushed to back of list successfully"
}
```

#### Pop a value from the front of a list

```
DELETE /api/v1/list/{key}/front
```

**cURL Example:**
```bash
curl -X DELETE http://localhost:8090/api/v1/list/mylist/front
```

**Response:**
```json
{
  "data": {
    "key": "mylist",
    "value": "first item"
  },
  "msg": "Value popped from front of list successfully"
}
```

#### Pop a value from the back of a list

```
DELETE /api/v1/list/{key}/back
```

**cURL Example:**
```bash
curl -X DELETE http://localhost:8090/api/v1/list/mylist/back
```

**Response:**
```json
{
  "data": {
    "key": "mylist",
    "value": "last item"
  },
  "msg": "Value popped from back of list successfully"
}
```

#### Get a range of values from a list

```
GET /api/v1/list/{key}/range?start={start}&end={end}
```

**cURL Example:**
```bash
curl -X GET "http://localhost:8090/api/v1/list/mylist/range?start=0&end=-1"
```

**Response:**
```json
{
  "data": {
    "key": "mylist",
    "start": 0,
    "end": -1,
    "values": ["first item", "middle item", "last item"]
  },
  "msg": "List range retrieved successfully"
}
```

### TTL Operations API

#### Set TTL for a key

```
PUT /api/v1/ttl/{key}
```

**Request Body:**
```json
{
  "ttl": 60
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8090/api/v1/ttl/greeting \
  -H "Content-Type: application/json" \
  -d '{"ttl": 60}'
```

**Response:**
```json
{
  "data": {
    "key": "greeting",
    "ttl": 60
  },
  "msg": "TTL set successfully"
}
```

#### Get TTL for a key

```
GET /api/v1/ttl/{key}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8090/api/v1/ttl/greeting
```

**Response:**
```json
{
  "data": {
    "key": "greeting",
    "ttl": 58.5
  },
  "msg": "TTL retrieved successfully"
}
```

#### Remove TTL for a key

```
DELETE /api/v1/ttl/{key}
```

**cURL Example:**
```bash
curl -X DELETE http://localhost:8090/api/v1/ttl/greeting
```

**Response:**
```json
{
  "data": {
    "key": "greeting"
  },
  "msg": "TTL removed successfully"
}
```

### General Operations API

#### Remove a key

```
DELETE /api/v1/key/{key}
```

**cURL Example:**
```bash
curl -X DELETE http://localhost:8090/api/v1/key/greeting
```

**Response:**
```json
{
  "data": {
    "key": "greeting"
  },
  "msg": "Key removed successfully"
}
```

#### Check if a key exists

```
GET /api/v1/key/{key}/exists
```

**cURL Example:**
```bash
curl -X GET http://localhost:8090/api/v1/key/greeting/exists
```

**Response:**
```json
{
  "data": {
    "key": "greeting",
    "exists": true
  },
  "msg": "Key existence checked"
}
```

#### Get the type of a key

```
GET /api/v1/key/{key}/type
```

**cURL Example:**
```bash
curl -X GET http://localhost:8090/api/v1/key/greeting/type
```

**Response:**
```json
{
  "data": {
    "key": "greeting",
    "type": "string"
  },
  "msg": "Key type retrieved successfully"
}
```

#### Clear all keys

```
DELETE /api/v1/keys
```

**cURL Example:**
```bash
curl -X DELETE http://localhost:8090/api/v1/keys
```

**Response:**
```json
{
  "data": {},
  "msg": "Cache cleared successfully"
}
```

## Running Locally with Docker 🐳

Gredis can be easily run locally using Docker. There are two main ways to run the application:

### Using Docker Directly

To build and run the application using Docker directly:

```bash
# Build the Docker image
docker build -t gredis .

# Run the container in development mode
docker run -p 8090:8090 -v $(pwd)/config.yaml:/app/config.yaml --name gredis-dev gredis

# Or run in production mode
docker run -p 8090:8090 -v $(pwd)/config.yaml:/app/config.yaml --name gredis-prod --target production gredis
```

### Using Docker Compose

For a more convenient setup, you can use Docker Compose:

```bash
# Start the application with Docker Compose
docker compose up --build

# Run in detached mode
docker compose up -d

# Stop and remove containers
docker compose down --remove-orphans --volumes
```

### Using Taskfile

If you have [Task](https://taskfile.dev/) installed, you can use the following commands to manage Docker Compose:

```bash
# Start the application with Docker Compose
task compose-up

# Stop and remove Docker Compose containers
task compose-down
```

Once the application is running, you can access the API at `http://localhost:8090`.

## Taskfile Commands 🛠️

Gredis uses [Task](https://taskfile.dev/) for common development operations. Below are the available commands:

### Lint

Run the linter to check code quality:

```bash
task lint
```

### Test

Run all tests with race detection enabled:

```bash
task test
```

### Docker Compose

Start the application with Docker Compose:

```bash
task compose-up
```

Stop and remove Docker Compose containers:

```bash
task compose-down
```

## Examples 📚

See the [examples](examples) directory for more detailed examples.

## License 📄

This project is licensed under the MIT License - see the LICENSE file for details.
