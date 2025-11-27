package cache

import (
	"context"
	"time"
)

// Cache interface defines the methods for caching operations
type Cache interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) (string, error)
	
	// Set sets a value in cache with expiration
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	
	// SetNX sets a value in cache only if it doesn't already exist
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	
	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error
	
	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)
	
	// FlushDB clears all keys from the current database
	FlushDB(ctx context.Context) error
	
	// Keys retrieves all keys matching the pattern
	Keys(ctx context.Context, pattern string) ([]string, error)
	
	// GetJSON retrieves a value from cache and unmarshals it to the provided struct
	GetJSON(ctx context.Context, key string, dest interface{}) error
	
	// SetJSON sets a value in cache with expiration after marshaling to JSON
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}
