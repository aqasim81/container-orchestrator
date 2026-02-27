package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// RunStoreTests runs a standard suite of tests against any Store implementation.
func RunStoreTests(t *testing.T, newStore func(t *testing.T) Store) {
	t.Helper()

	t.Run("PutAndGet", func(t *testing.T) {
		s := newStore(t)
		defer func() { require.NoError(t, s.Close()) }()

		err := s.Put("test", "key1", []byte("value1"))
		require.NoError(t, err)

		val, err := s.Get("test", "key1")
		require.NoError(t, err)
		assert.Equal(t, []byte("value1"), val)
	})

	t.Run("GetFromNonexistentBucket", func(t *testing.T) {
		s := newStore(t)
		defer func() { require.NoError(t, s.Close()) }()

		_, err := s.Get("nonexistent", "key1")
		assert.ErrorIs(t, err, ErrBucketNotFound)
	})

	t.Run("GetNonexistentKey", func(t *testing.T) {
		s := newStore(t)
		defer func() { require.NoError(t, s.Close()) }()

		err := s.Put("test", "key1", []byte("value1"))
		require.NoError(t, err)

		_, err = s.Get("test", "missing")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("PutOverwrite", func(t *testing.T) {
		s := newStore(t)
		defer func() { require.NoError(t, s.Close()) }()

		require.NoError(t, s.Put("test", "key1", []byte("v1")))
		require.NoError(t, s.Put("test", "key1", []byte("v2")))

		val, err := s.Get("test", "key1")
		require.NoError(t, err)
		assert.Equal(t, []byte("v2"), val)
	})

	t.Run("Delete", func(t *testing.T) {
		s := newStore(t)
		defer func() { require.NoError(t, s.Close()) }()

		require.NoError(t, s.Put("test", "key1", []byte("value1")))

		err := s.Delete("test", "key1")
		require.NoError(t, err)

		_, err = s.Get("test", "key1")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("DeleteFromNonexistentBucket", func(t *testing.T) {
		s := newStore(t)
		defer func() { require.NoError(t, s.Close()) }()

		err := s.Delete("nonexistent", "key1")
		assert.ErrorIs(t, err, ErrBucketNotFound)
	})

	t.Run("DeleteNonexistentKey", func(t *testing.T) {
		s := newStore(t)
		defer func() { require.NoError(t, s.Close()) }()

		require.NoError(t, s.Put("test", "key1", []byte("value1")))

		err := s.Delete("test", "missing")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("ListAll", func(t *testing.T) {
		s := newStore(t)
		defer func() { require.NoError(t, s.Close()) }()

		require.NoError(t, s.Put("test", "a:1", []byte("v1")))
		require.NoError(t, s.Put("test", "a:2", []byte("v2")))
		require.NoError(t, s.Put("test", "b:1", []byte("v3")))

		results, err := s.List("test", "")
		require.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("ListWithPrefix", func(t *testing.T) {
		s := newStore(t)
		defer func() { require.NoError(t, s.Close()) }()

		require.NoError(t, s.Put("test", "container:1", []byte("c1")))
		require.NoError(t, s.Put("test", "container:2", []byte("c2")))
		require.NoError(t, s.Put("test", "node:1", []byte("n1")))

		results, err := s.List("test", "container:")
		require.NoError(t, err)
		assert.Len(t, results, 2)

		for _, kv := range results {
			assert.Contains(t, kv.Key, "container:")
		}
	})

	t.Run("ListNonexistentBucket", func(t *testing.T) {
		s := newStore(t)
		defer func() { require.NoError(t, s.Close()) }()

		_, err := s.List("nonexistent", "")
		assert.ErrorIs(t, err, ErrBucketNotFound)
	})

	t.Run("MultipleBuckets", func(t *testing.T) {
		s := newStore(t)
		defer func() { require.NoError(t, s.Close()) }()

		require.NoError(t, s.Put("nodes", "n1", []byte("node1")))
		require.NoError(t, s.Put("containers", "c1", []byte("container1")))

		val, err := s.Get("nodes", "n1")
		require.NoError(t, err)
		assert.Equal(t, []byte("node1"), val)

		val, err = s.Get("containers", "c1")
		require.NoError(t, err)
		assert.Equal(t, []byte("container1"), val)

		_, err = s.Get("nodes", "c1")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("ValueIsolation", func(t *testing.T) {
		s := newStore(t)
		defer func() { require.NoError(t, s.Close()) }()

		original := []byte("original")
		require.NoError(t, s.Put("test", "key1", original))

		// Mutate original slice — should not affect stored value.
		original[0] = 'X'

		val, err := s.Get("test", "key1")
		require.NoError(t, err)
		assert.Equal(t, []byte("original"), val)

		// Mutate returned slice — should not affect stored value.
		val[0] = 'Y'

		val2, err := s.Get("test", "key1")
		require.NoError(t, err)
		assert.Equal(t, []byte("original"), val2)
	})
}
