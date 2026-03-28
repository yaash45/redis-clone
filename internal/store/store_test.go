package store

import (
	"testing"
)

func TestEmptyStoreAccess(t *testing.T) {
	kvs := NewStore()

	val := kvs.Get("foo")

	if val != "" {
		t.Error("Store is not empty.")
	}
}

func TestBasicApi(t *testing.T) {
	kvs := NewStore()

	kvs.Set("foo", "bar")

	if kvs.Get("foo") != "bar" {
		t.Fail()
	}

	kvs.Set("hi", "hello")

	if kvs.Get("hi") != "hello" {
		t.Fail()
	}

	// Overwrite a key's value
	kvs.Set("hi", "bye")

	if kvs.Get("hi") != "bye" {
		t.Fail()
	}

	// Delete a key
	kvs.Delete("foo")

	if kvs.Get("foo") != "" {
		t.Fail()
	}
}
