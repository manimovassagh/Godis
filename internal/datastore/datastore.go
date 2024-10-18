package datastore

import "sync"

type DataStore struct {
	data map[string]string
	mu   sync.RWMutex
}

var (
	instance *DataStore
	once     sync.Once
)

// GetDataStore returns a singleton instance of the DataStore. It is safe to call
// from multiple goroutines.
func GetDataStore() *DataStore {
	once.Do(func() {
		instance = &DataStore{
			data: make(map[string]string),
		}
	})
	return instance
}

// Set sets the given key-value pair in the in-memory data store. It is thread-safe
// and can be safely called from multiple goroutines concurrently.
func (ds *DataStore) Set(key, value string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.data[key] = value
}

// Get looks up the given key in the in-memory data store and returns the associated
// value, or ("", false) if the key is not found. It is thread-safe and can be
// safely called from multiple goroutines concurrently.
func (ds *DataStore) Get(key string) (string, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	value, found := ds.data[key]
	return value, found
}
