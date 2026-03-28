package store

// Represents an in-memory Key-Value store
type KVStore struct {
	store map[string]string
}

// Creates a new instance of a Key-Value Store.
func NewStore() *KVStore {
	return &KVStore{store: make(map[string]string)}
}

// Get the value for a given key from the Key-Value Store.
func (kvs *KVStore) Get(key string) string {
	return kvs.store[key]
}

// Stores a Key-Value pair in the store, overriding
// an existing value.
func (kvs *KVStore) Set(key string, value string) {
	kvs.store[key] = value
}

// Deletes a Key-Value pair from the store
func (kvs *KVStore) Delete(key string) {
	delete(kvs.store, key)
}
