// Package store provides a key-value storage interface and implementations.
package store

import "errors"

// Sentinel errors for store operations.
var (
	ErrNotFound       = errors.New("key not found")
	ErrBucketNotFound = errors.New("bucket not found")
)

// KV represents a key-value pair returned from list operations.
type KV struct {
	Key   string
	Value []byte
}

// Store defines the interface for key-value storage operations.
// All values are stored as raw bytes; callers handle serialization.
type Store interface {
	// Get retrieves a value by bucket and key.
	// Returns ErrBucketNotFound if the bucket does not exist.
	// Returns ErrNotFound if the key does not exist within the bucket.
	Get(bucket, key string) ([]byte, error)

	// Put stores a value in the given bucket under the given key.
	// The bucket is created automatically if it does not exist.
	Put(bucket, key string, value []byte) error

	// Delete removes a key from a bucket.
	// Returns ErrBucketNotFound if the bucket does not exist.
	// Returns ErrNotFound if the key does not exist.
	Delete(bucket, key string) error

	// List returns all key-value pairs in a bucket whose keys start with prefix.
	// An empty prefix returns all entries in the bucket.
	// Returns ErrBucketNotFound if the bucket does not exist.
	List(bucket, prefix string) ([]KV, error)

	// Close releases any resources held by the store.
	Close() error
}
