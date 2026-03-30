// Contains the core key-value store implementation
package store

import (
	"errors"

	"github.com/yaash45/redis/internal/command"
)

// Represents an in-memory Key-Value store
type KVStore struct {
	store map[string]string
}

// Creates a new instance of a Key-Value Store.
func NewStore() *KVStore {
	return &KVStore{store: make(map[string]string)}
}

// Get the value for a given key from the Key-Value Store.
func (kvs *KVStore) Get(key string) (string, error) {
	val := kvs.store[key]

	if val == "" {
		return "", errors.New("Key not found.")
	}

	return kvs.store[key], nil
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

// Processes the input commands and executes the relevant operation
// on the internal in-memory key-value store.
func (kvs *KVStore) ProcessCmd(cmd *command.Command) (string, error) {
	name := cmd.Name()
	key := cmd.Arg1()
	value := cmd.Arg2()

	switch name {
	case "GET":
		lookup_result, err := kvs.Get(key)
		if err != nil {
			return "", err
		} else {
			return lookup_result, nil
		}
	case "SET":
		if key == "" || value == "" {
			return "", errors.New("[error] 'SET' command requires non-empty key and value.")
		} else {
			kvs.Set(key, value)
			return "[ok]", nil
		}
	case "DEL":
		kvs.Delete(key)
		return "ok", nil
	default:
		return "", errors.New("Unrecognized command. Available commands: SET/GET/DEL <KEY> <VALUE>.")
	}
}
