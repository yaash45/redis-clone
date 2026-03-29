package store

import (
	"testing"
)

func TestEmptyStoreAccess(t *testing.T) {
	kvs := NewStore()

	_, err := kvs.Get("foo")

	if err == nil {
		t.Error("Store is not empty.")
	}
}

func TestBasicApi(t *testing.T) {
	kvs := NewStore()

	kvs.Set("foo", "bar")
	val, err := kvs.Get("foo")

	if err != nil && val != "bar" {
		t.Fail()
	}

	kvs.Set("hi", "hello")
	val, err = kvs.Get("hi")

	if err != nil && val != "hello" {
		t.Fail()
	}

	// Overwrite a key's value
	kvs.Set("hi", "bye")
	val, err = kvs.Get("hi")

	if err != nil && val != "bye" {
		t.Fail()
	}

	// Delete a key
	kvs.Delete("foo")
	_, err = kvs.Get("foo")

	if err == nil {
		t.Fail()
	}
}
