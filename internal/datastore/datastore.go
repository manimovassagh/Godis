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

func GetDataStore() *DataStore {
    once.Do(func() {
        instance = &DataStore{
            data: make(map[string]string),
        }
    })
    return instance
}

func (ds *DataStore) Set(key, value string) {
    ds.mu.Lock()
    defer ds.mu.Unlock()
    ds.data[key] = value
}

func (ds *DataStore) Get(key string) (string, bool) {
    ds.mu.RLock()
    defer ds.mu.RUnlock()
    value, found := ds.data[key]
    return value, found
}