package store

import (
	"strings"
	"sync"
)

// MemoryStore implements Store using an in-memory map. Intended for testing.
type MemoryStore struct {
	buckets map[string]map[string][]byte
	mu      sync.RWMutex
}

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		buckets: make(map[string]map[string][]byte),
	}
}

// Get retrieves a value by bucket and key.
func (s *MemoryStore) Get(bucket, key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.buckets[bucket]
	if !ok {
		return nil, ErrBucketNotFound
	}

	v, ok := b[key]
	if !ok {
		return nil, ErrNotFound
	}

	cp := make([]byte, len(v))
	copy(cp, v)
	return cp, nil
}

// Put stores a value in the given bucket under the given key.
func (s *MemoryStore) Put(bucket, key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.buckets[bucket]
	if !ok {
		b = make(map[string][]byte)
		s.buckets[bucket] = b
	}

	cp := make([]byte, len(value))
	copy(cp, value)
	b[key] = cp
	return nil
}

// Delete removes a key from a bucket.
func (s *MemoryStore) Delete(bucket, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.buckets[bucket]
	if !ok {
		return ErrBucketNotFound
	}

	if _, ok := b[key]; !ok {
		return ErrNotFound
	}

	delete(b, key)
	return nil
}

// List returns all key-value pairs in a bucket whose keys start with prefix.
func (s *MemoryStore) List(bucket, prefix string) ([]KV, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.buckets[bucket]
	if !ok {
		return nil, ErrBucketNotFound
	}

	var results []KV
	for k, v := range b {
		if prefix != "" && !strings.HasPrefix(k, prefix) {
			continue
		}
		cp := make([]byte, len(v))
		copy(cp, v)
		results = append(results, KV{Key: k, Value: cp})
	}

	return results, nil
}

// Close is a no-op for the in-memory store.
func (s *MemoryStore) Close() error {
	return nil
}
