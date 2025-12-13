package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// FallbackCache implements the Cache interface using in-memory storage
// This is used when Redis is unavailable
type FallbackCache struct {
	data     map[string]*cacheItem
	mutex    sync.RWMutex
	cleanup  chan struct{}
}

type cacheItem struct {
	value     string
	expiry    time.Time
	isExpired bool
}

// NewFallbackCache creates a new instance of FallbackCache
func NewFallbackCache() *FallbackCache {
	fc := &FallbackCache{
		data:    make(map[string]*cacheItem),
		cleanup: make(chan struct{}),
	}
	
	// Start cleanup goroutine to remove expired items
	go fc.startCleanup()
	
	return fc
}

// startCleanup periodically removes expired cache items
func (fc *FallbackCache) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute) // Clean up every 5 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fc.mutex.Lock()
			now := time.Now()
			for key, item := range fc.data {
				if now.After(item.expiry) {
					delete(fc.data, key)
				}
			}
			fc.mutex.Unlock()
		case <-fc.cleanup:
			return // Stop cleanup when cleanup channel is closed
		}
	}
}

// Close stops the cleanup goroutine
func (fc *FallbackCache) Close() {
	close(fc.cleanup)
}

// Get retrieves a value from cache
func (fc *FallbackCache) Get(ctx context.Context, key string) (string, error) {
	fc.mutex.RLock()
	defer fc.mutex.RUnlock()

	item, exists := fc.data[key]
	if !exists {
		return "", fmt.Errorf("key %s not found in cache", key)
	}

	if time.Now().After(item.expiry) {
		// Mark as expired but don't delete immediately to avoid race conditions
		// Cleanup will remove it later
		return "", fmt.Errorf("key %s has expired", key)
	}

	return item.value, nil
}

// Set sets a value in cache with expiration
func (fc *FallbackCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	var valueStr string
	switch v := value.(type) {
	case string:
		valueStr = v
	default:
		// Convert to string representation
		bytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %v", err)
		}
		valueStr = string(bytes)
	}

	fc.data[key] = &cacheItem{
		value:  valueStr,
		expiry: time.Now().Add(expiration),
	}

	return nil
}

// SetNX sets a value in cache only if it doesn't already exist
func (fc *FallbackCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	if _, exists := fc.data[key]; exists {
		return false, nil // Key already exists
	}

	var valueStr string
	switch v := value.(type) {
	case string:
		valueStr = v
	default:
		// Convert to string representation
		bytes, err := json.Marshal(v)
		if err != nil {
			return false, fmt.Errorf("failed to marshal value: %v", err)
		}
		valueStr = string(bytes)
	}

	fc.data[key] = &cacheItem{
		value:  valueStr,
		expiry: time.Now().Add(expiration),
	}

	return true, nil
}

// Delete removes a value from cache
func (fc *FallbackCache) Delete(ctx context.Context, key string) error {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	delete(fc.data, key)
	return nil
}

// Exists checks if a key exists in cache
func (fc *FallbackCache) Exists(ctx context.Context, key string) (bool, error) {
	fc.mutex.RLock()
	defer fc.mutex.RUnlock()

	item, exists := fc.data[key]
	if !exists {
		return false, nil
	}

	// Check if the item has expired
	if time.Now().After(item.expiry) {
		return false, nil
	}

	return true, nil
}

// FlushDB clears all keys from the cache
func (fc *FallbackCache) FlushDB(ctx context.Context) error {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	fc.data = make(map[string]*cacheItem)
	return nil
}

// Keys retrieves all keys matching the pattern
// Note: This is a simplified implementation that doesn't support complex patterns
func (fc *FallbackCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	fc.mutex.RLock()
	defer fc.mutex.RUnlock()

	var keys []string
	for key := range fc.data {
		// Simple pattern matching - only supports '*' at the end
		if matchesPattern(key, pattern) {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

// matchesPattern checks if a key matches a simple pattern
// Supports '*' wildcard at the end only
func matchesPattern(key, pattern string) bool {
	if pattern == "*" {
		return true
	}

	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(key) >= len(prefix) && key[:len(prefix)] == prefix
	}

	return key == pattern
}

// GetJSON retrieves a value from cache and unmarshals it to the provided struct
func (fc *FallbackCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	value, err := fc.Get(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(value), dest)
}

// SetJSON sets a value in cache with expiration after marshaling to JSON
func (fc *FallbackCache) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return fc.Set(ctx, key, value, expiration)
}