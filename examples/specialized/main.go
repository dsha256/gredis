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

	fmt.Println("=== String Operations with StringClient ===")
	// Set a string value
	err := strClient.Set("greeting", "Hello, World!")
	if err != nil {
		fmt.Printf("Error setting key: %v\n", err)
	}

	// Get the value
	value, err := strClient.Get("greeting")
	if err != nil {
		fmt.Printf("Error getting key: %v\n", err)
	} else {
		fmt.Printf("greeting = %s\n", value)
	}

	// Update the value
	err = strClient.Update("greeting", "Hello, Gredis!")
	if err != nil {
		fmt.Printf("Error updating key: %v\n", err)
	}

	// Get the updated value
	value, err = strClient.Get("greeting")
	if err != nil {
		fmt.Printf("Error getting key: %v\n", err)
	} else {
		fmt.Printf("greeting (updated) = %s\n", value)
	}

	fmt.Println("\n=== TTL Operations with TTLClient ===")
	// Set a value with TTL
	err = strClient.SetWithTTL("temp", "This will expire", 2*time.Second)
	if err != nil {
		fmt.Printf("Error setting key with TTL: %v\n", err)
	}

	// Get the TTL
	ttl, err := ttlClient.GetTTL("temp")
	if err != nil {
		fmt.Printf("Error getting TTL: %v\n", err)
	} else {
		fmt.Printf("TTL for 'temp': %v\n", ttl)
	}

	// Wait for the key to expire
	fmt.Println("Waiting for 'temp' to expire...")
	time.Sleep(3 * time.Second)

	// Try to get the expired value
	_, err = strClient.Get("temp")
	if err != nil {
		fmt.Printf("As expected, 'temp' has expired: %v\n", err)
	}

	fmt.Println("\n=== List Operations with ListClient ===")
	// Create a list
	err = listClient.PushBack("mylist", "first")
	if err != nil {
		fmt.Printf("Error pushing to list: %v\n", err)
	}

	err = listClient.PushBack("mylist", "second")
	if err != nil {
		fmt.Printf("Error pushing to list: %v\n", err)
	}

	err = listClient.PushFront("mylist", "zero")
	if err != nil {
		fmt.Printf("Error pushing to list: %v\n", err)
	}

	// Get the list range
	items, err := listClient.ListRange("mylist", 0, -1)
	if err != nil {
		fmt.Printf("Error getting list range: %v\n", err)
	} else {
		fmt.Printf("List items: %v\n", items)
	}

	// Pop from the list
	item, err := listClient.PopFront("mylist")
	if err != nil {
		fmt.Printf("Error popping from list: %v\n", err)
	} else {
		fmt.Printf("Popped from front: %s\n", item)
	}

	item, err = listClient.PopBack("mylist")
	if err != nil {
		fmt.Printf("Error popping from list: %v\n", err)
	} else {
		fmt.Printf("Popped from back: %s\n", item)
	}

	// Get the updated list
	items, err = listClient.ListRange("mylist", 0, -1)
	if err != nil {
		fmt.Printf("Error getting list range: %v\n", err)
	} else {
		fmt.Printf("List items after popping: %v\n", items)
	}

	fmt.Println("\n=== Type and Exists Operations ===")
	// Check if keys exist
	fmt.Printf("'greeting' exists: %v\n", c.Exists("greeting"))
	fmt.Printf("'nonexistent' exists: %v\n", c.Exists("nonexistent"))

	// Get the type of keys
	greetingType, err := c.Type("greeting")
	if err != nil {
		fmt.Printf("Error getting type: %v\n", err)
	} else {
		fmt.Printf("Type of 'greeting': %v\n", greetingType)
	}

	listType, err := c.Type("mylist")
	if err != nil {
		fmt.Printf("Error getting type: %v\n", err)
	} else {
		fmt.Printf("Type of 'mylist': %v\n", listType)
	}

	fmt.Println("\n=== Cleanup ===")
	// Remove a key
	err = c.Remove("greeting")
	if err != nil {
		fmt.Printf("Error removing key: %v\n", err)
	}

	// Check if the key still exists
	fmt.Printf("'greeting' exists after removal: %v\n", c.Exists("greeting"))

	// Clear all keys
	err = c.Clear()
	if err != nil {
		fmt.Printf("Error clearing cache: %v\n", err)
	}

	// Check if any keys still exist
	fmt.Printf("'mylist' exists after clear: %v\n", c.Exists("mylist"))
}
