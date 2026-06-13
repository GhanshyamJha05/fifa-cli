package repository

import (
	"time"

	"github.com/GhanshyamJha05/fifa-cli/internal/cache"
	"github.com/GhanshyamJha05/fifa-cli/internal/provider"
)

// CacheRepository wraps BoltDB cache for dependency injection.
type CacheRepository struct {
	*cache.Store
}

// OpenCache opens the cache store at path with the given TTL.
func OpenCache(path string, ttl time.Duration) (*CacheRepository, error) {
	store, err := cache.Open(path, ttl)
	if err != nil {
		return nil, err
	}
	return &CacheRepository{Store: store}, nil
}

// Close closes the underlying store.
func (r *CacheRepository) Close() error {
	return r.Store.Close()
}

// Ensure CacheRepository implements provider.Cacher.
var _ provider.Cacher = (*CacheRepository)(nil)
