package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newBoltTestStore(t *testing.T) Store {
	t.Helper()
	dir := t.TempDir()
	s, err := NewBoltStore(dir)
	require.NoError(t, err)
	return s
}

func TestBoltStore(t *testing.T) {
	RunStoreTests(t, newBoltTestStore)
}

func TestBoltStore_Persistence(t *testing.T) {
	dir := t.TempDir()

	// Write data and close.
	s1, err := NewBoltStore(dir)
	require.NoError(t, err)
	require.NoError(t, s1.Put("test", "key1", []byte("persisted")))
	require.NoError(t, s1.Close())

	// Reopen and verify data survives.
	s2, err := NewBoltStore(dir)
	require.NoError(t, err)
	defer func() { require.NoError(t, s2.Close()) }()

	val, err := s2.Get("test", "key1")
	require.NoError(t, err)
	assert.Equal(t, []byte("persisted"), val)
}
