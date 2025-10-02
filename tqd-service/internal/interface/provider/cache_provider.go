package provider

import (
	"context"
	"time"
)

// CacheProvider defines the interface for cache operations
type CacheProvider interface {
	// Set stores a value in cache with expiration
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Get retrieves a value from cache
	Get(ctx context.Context, key string, dest interface{}) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)

	// SetNX sets a key only if it doesn't exist
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)

	// Increment increments a numeric value in cache
	Increment(ctx context.Context, key string, delta int64) (int64, error)

	// Expire sets expiration time for a key
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// GetTTL gets time to live for a key
	GetTTL(ctx context.Context, key string) (time.Duration, error)
}
