package cache

import (
	"encoding/json"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

var bucket = []byte("fifa")

// Store provides BoltDB-backed caching.
type Store struct {
	db  *bolt.DB
	ttl time.Duration
}

type entry struct {
	Data      json.RawMessage `json:"data"`
	FetchedAt time.Time       `json:"fetched_at"`
}

// Open creates or opens a cache database.
func Open(path string, ttl time.Duration) (*Store, error) {
	db, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("open bolt db: %w", err)
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	}); err != nil {
		db.Close()
		return nil, err
	}

	return &Store{db: db, ttl: ttl}, nil
}

// Close closes the database.
func (s *Store) Close() error {
	return s.db.Close()
}

// Get retrieves cached data if fresh.
func (s *Store) Get(key string, dest any) (bool, error) {
	var raw json.RawMessage
	found := false

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		v := b.Get([]byte(key))
		if v == nil {
			return nil
		}
		var e entry
		if err := json.Unmarshal(v, &e); err != nil {
			return err
		}
		if time.Since(e.FetchedAt) > s.ttl {
			return nil
		}
		raw = e.Data
		found = true
		return nil
	})
	if err != nil || !found {
		return false, err
	}

	if err := json.Unmarshal(raw, dest); err != nil {
		return false, err
	}
	return true, nil
}

// GetStale retrieves cached data regardless of TTL (offline mode).
func (s *Store) GetStale(key string, dest any) (bool, error) {
	var raw json.RawMessage
	found := false

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		v := b.Get([]byte(key))
		if v == nil {
			return nil
		}
		var e entry
		if err := json.Unmarshal(v, &e); err != nil {
			return err
		}
		raw = e.Data
		found = true
		return nil
	})
	if err != nil || !found {
		return false, err
	}

	return true, json.Unmarshal(raw, dest)
}

// Set stores data in cache.
func (s *Store) Set(key string, data any) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	e := entry{Data: raw, FetchedAt: time.Now()}
	encoded, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.Put([]byte(key), encoded)
	})
}
