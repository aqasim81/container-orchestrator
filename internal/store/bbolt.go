package store

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

// BoltStore implements Store using bbolt as the backing engine.
type BoltStore struct {
	db *bolt.DB
}

// NewBoltStore opens or creates a bbolt database at the given directory.
// The directory is created if it does not exist.
func NewBoltStore(dataDir string) (*BoltStore, error) {
	if err := os.MkdirAll(dataDir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create data directory %s: %w", dataDir, err)
	}

	dbPath := filepath.Join(dataDir, "orchestrator.db")
	db, err := bolt.Open(dbPath, 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database %s: %w", dbPath, err)
	}

	return &BoltStore{db: db}, nil
}

// Get retrieves a value by bucket and key.
func (s *BoltStore) Get(bucket, key string) ([]byte, error) {
	var value []byte

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrBucketNotFound
		}

		v := b.Get([]byte(key))
		if v == nil {
			return ErrNotFound
		}

		value = make([]byte, len(v))
		copy(value, v)
		return nil
	})

	return value, err
}

// Put stores a value in the given bucket under the given key.
func (s *BoltStore) Put(bucket, key string, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
		}
		return b.Put([]byte(key), value)
	})
}

// Delete removes a key from a bucket.
func (s *BoltStore) Delete(bucket, key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrBucketNotFound
		}

		v := b.Get([]byte(key))
		if v == nil {
			return ErrNotFound
		}

		return b.Delete([]byte(key))
	})
}

// List returns all key-value pairs in a bucket whose keys start with prefix.
func (s *BoltStore) List(bucket, prefix string) ([]KV, error) {
	var results []KV

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrBucketNotFound
		}

		c := b.Cursor()

		if prefix == "" {
			for k, v := c.First(); k != nil; k, v = c.Next() {
				val := make([]byte, len(v))
				copy(val, v)
				results = append(results, KV{
					Key:   string(k),
					Value: val,
				})
			}
		} else {
			prefixBytes := []byte(prefix)
			for k, v := c.Seek(prefixBytes); k != nil && bytes.HasPrefix(k, prefixBytes); k, v = c.Next() {
				val := make([]byte, len(v))
				copy(val, v)
				results = append(results, KV{
					Key:   string(k),
					Value: val,
				})
			}
		}

		return nil
	})

	return results, err
}

// Close closes the underlying bbolt database.
func (s *BoltStore) Close() error {
	return s.db.Close()
}
