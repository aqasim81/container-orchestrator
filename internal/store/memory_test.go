package store

import "testing"

func newMemoryTestStore(t *testing.T) Store {
	t.Helper()
	return NewMemoryStore()
}

func TestMemoryStore(t *testing.T) {
	RunStoreTests(t, newMemoryTestStore)
}
