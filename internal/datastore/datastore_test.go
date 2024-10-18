package datastore

import (
	"sync"
	"testing"
)

// TestGetDataStore tests the singleton behavior of GetDataStore
func TestGetDataStore(t *testing.T) {
	ds1 := GetDataStore()
	ds2 := GetDataStore()

	if ds1 != ds2 {
		t.Errorf("Expected both instances to be the same, got different instances")
	}
}

// TestSetAndGet tests the Set and Get methods
func TestSetAndGet(t *testing.T) {
	ds := GetDataStore()

	// Set key-value pair
	ds.Set("foo", "bar")

	// Test Get for existing key
	value, found := ds.Get("foo")
	if !found {
		t.Errorf("Expected to find key 'foo', but it was not found")
	}
	if value != "bar" {
		t.Errorf("Expected value 'bar', got %s", value)
	}

	// Test Get for non-existing key
	_, found = ds.Get("nonexistent")
	if found {
		t.Errorf("Expected 'nonexistent' key not to be found, but it was found")
	}
}

// TestConcurrentAccess tests concurrent access to the DataStore
func TestConcurrentAccess(t *testing.T) {
	ds := GetDataStore()
	var wg sync.WaitGroup
	const numGoroutines = 100

	// Launch multiple goroutines to set values
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ds.Set(string(rune('A'+i)), string(rune('a'+i)))
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Test if the values were set correctly
	for i := 0; i < numGoroutines; i++ {
		key := string(rune('A' + i))
		value, found := ds.Get(key)
		if !found {
			t.Errorf("Expected to find key '%s', but it was not found", key)
		}
		expectedValue := string(rune('a' + i))
		if value != expectedValue {
			t.Errorf("Expected value '%s', got '%s'", expectedValue, value)
		}
	}
}