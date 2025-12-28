package cache

import (
	"context"
	"log"
	"sync"
	"time"
)

// AdaptiveCache wraps Redis and Fallback cache with automatic failover
type AdaptiveCache struct {
	redisCache    *RedisCache
	fallbackCache *FallbackCache
	mutex         sync.RWMutex
	usingFallback bool
}

// NewAdaptiveCache creates a new adaptive cache that switches between Redis and in-memory based on Redis availability
func NewAdaptiveCache(redisManager *RedisManager) *AdaptiveCache {
	var redisCache *RedisCache
	if redisManager != nil {
		redisCache = NewRedisCache(redisManager)
	} else {
		redisCache = nil // Will be nil if Redis manager is not available
	}
	fallbackCache := NewFallbackCache()

	ac := &AdaptiveCache{
		redisCache:    redisCache,
		fallbackCache: fallbackCache,
		usingFallback: true, // Start with fallback if Redis manager is nil
	}

	// Start monitoring Redis health to automatically switch between implementations
	// Only if Redis manager is available
	if redisManager != nil {
		go ac.monitorRedisHealth(redisManager)
	}

	return ac
}

// monitorRedisHealth continuously monitors Redis health
func (ac *AdaptiveCache) monitorRedisHealth(redisManager *RedisManager) {
	ticker := time.NewTicker(10 * time.Second) // Check every 10 seconds
	defer ticker.Stop()

	for {
		time.Sleep(10 * time.Second)

		isHealthy := redisManager.IsHealthy()

		ac.mutex.Lock()
		wasUsingFallback := ac.usingFallback
		shouldUseFallback := !isHealthy

		if wasUsingFallback != shouldUseFallback {
			ac.usingFallback = shouldUseFallback
			if shouldUseFallback {
				log.Println("Switching to fallback cache due to Redis unavailability")
			} else {
				log.Println("Switching back to Redis cache")
			}
		}
		ac.mutex.Unlock()
	}
}

// Get retrieves a value from cache
func (ac *AdaptiveCache) Get(ctx context.Context, key string) (string, error) {
	ac.mutex.RLock()
	usingFallback := ac.usingFallback
	ac.mutex.RUnlock()

	if usingFallback || ac.redisCache == nil {
		return ac.fallbackCache.Get(ctx, key)
	}

	value, err := ac.redisCache.Get(ctx, key)
	if err != nil {
		// If Redis fails, try fallback cache
		log.Printf("Redis error, trying fallback: %v", err)
		ac.mutex.Lock()
		ac.usingFallback = true
		ac.mutex.Unlock()

		return ac.fallbackCache.Get(ctx, key)
	}

	return value, err
}

// Set sets a value in cache with expiration
func (ac *AdaptiveCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ac.mutex.RLock()
	usingFallback := ac.usingFallback
	ac.mutex.RUnlock()

	var primaryErr error
	var fallbackErr error

	// Try primary cache first if it's available and we're not using fallback
	if !usingFallback && ac.redisCache != nil {
		primaryErr = ac.redisCache.Set(ctx, key, value, expiration)
		if primaryErr != nil {
			// Switch to fallback on error
			log.Printf("Redis error during Set, switching to fallback: %v", primaryErr)
			ac.mutex.Lock()
			ac.usingFallback = true
			ac.mutex.Unlock()
		} else {
			// Set in fallback cache as well for consistency when Redis is working
			ac.fallbackCache.Set(ctx, key, value, expiration)
			return nil // Success with primary cache
		}
	}

	// Set in fallback cache as well for consistency
	fallbackErr = ac.fallbackCache.Set(ctx, key, value, expiration)
	if fallbackErr != nil {
		log.Printf("Fallback cache error during Set: %v", fallbackErr)
	}

	// Return the primary error if we were trying to use Redis and it failed
	if primaryErr != nil {
		return primaryErr
	}

	return fallbackErr
}

// SetNX sets a value in cache only if it doesn't already exist
func (ac *AdaptiveCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	ac.mutex.RLock()
	usingFallback := ac.usingFallback
	ac.mutex.RUnlock()

	if usingFallback || ac.redisCache == nil {
		return ac.fallbackCache.SetNX(ctx, key, value, expiration)
	}

	success, err := ac.redisCache.SetNX(ctx, key, value, expiration)
	if err != nil {
		// If Redis fails, try fallback cache
		log.Printf("Redis error during SetNX, trying fallback: %v", err)
		ac.mutex.Lock()
		ac.usingFallback = true
		ac.mutex.Unlock()

		return ac.fallbackCache.SetNX(ctx, key, value, expiration)
	}

	return success, err
}

// Delete removes a value from cache
func (ac *AdaptiveCache) Delete(ctx context.Context, key string) error {
	var primaryErr error
	var fallbackErr error

	// Always delete from primary cache if available (regardless of fallback status)
	if ac.redisCache != nil {
		primaryErr = ac.redisCache.Delete(ctx, key)
		if primaryErr != nil {
			// Log the error but continue to delete from fallback
			log.Printf("Redis error during Delete: %v", primaryErr)
		}
	}

	// Always delete from fallback cache as well to ensure complete deletion
	fallbackErr = ac.fallbackCache.Delete(ctx, key)
	if fallbackErr != nil {
		log.Printf("Fallback cache error during Delete: %v", fallbackErr)
	}

	// Return the primary error if there was one (prioritizing Redis errors)
	if primaryErr != nil {
		return primaryErr
	}

	return fallbackErr
}

// Exists checks if a key exists in cache
func (ac *AdaptiveCache) Exists(ctx context.Context, key string) (bool, error) {
	ac.mutex.RLock()
	usingFallback := ac.usingFallback
	ac.mutex.RUnlock()

	if usingFallback || ac.redisCache == nil {
		return ac.fallbackCache.Exists(ctx, key)
	}

	exists, err := ac.redisCache.Exists(ctx, key)
	if err != nil {
		// If Redis fails, try fallback cache
		log.Printf("Redis error during Exists, trying fallback: %v", err)
		ac.mutex.Lock()
		ac.usingFallback = true
		ac.mutex.Unlock()

		return ac.fallbackCache.Exists(ctx, key)
	}

	return exists, err
}

// FlushDB clears all keys from the current database
func (ac *AdaptiveCache) FlushDB(ctx context.Context) error {
	var primaryErr error
	var fallbackErr error

	// Always flush primary cache if available
	if ac.redisCache != nil {
		primaryErr = ac.redisCache.FlushDB(ctx)
		if primaryErr != nil {
			log.Printf("Redis error during FlushDB: %v", primaryErr)
		}
	}

	// Always flush fallback cache as well to ensure complete flush
	fallbackErr = ac.fallbackCache.FlushDB(ctx)
	if fallbackErr != nil {
		log.Printf("Fallback cache error during FlushDB: %v", fallbackErr)
	}

	// Return the primary error if there was one (prioritizing Redis errors)
	if primaryErr != nil {
		return primaryErr
	}

	return fallbackErr
}

// Keys retrieves all keys matching the pattern
// This implementation tries both caches to ensure comprehensive key retrieval
func (ac *AdaptiveCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	ac.mutex.RLock()
	usingFallback := ac.usingFallback
	ac.mutex.RUnlock()

	var allKeys []string
	var primaryErr error

	// Try primary cache if available
	if !usingFallback && ac.redisCache != nil {
		keys, err := ac.redisCache.Keys(ctx, pattern)
		if err != nil {
			// Log the error but continue to try fallback
			log.Printf("Redis error during Keys: %v", err)
			primaryErr = err
		} else {
			allKeys = append(allKeys, keys...)
		}
	}

	// Always try fallback cache as well to ensure comprehensive results
	fallbackKeys, fallbackErr := ac.fallbackCache.Keys(ctx, pattern)
	if fallbackErr != nil {
		log.Printf("Fallback cache error during Keys: %v", fallbackErr)
	} else {
		allKeys = append(allKeys, fallbackKeys...)
	}

	// If we have results from either cache, return them
	if len(allKeys) > 0 {
		// Remove duplicates
		uniqueKeys := make([]string, 0, len(allKeys))
		seen := make(map[string]bool)
		for _, key := range allKeys {
			if !seen[key] {
				uniqueKeys = append(uniqueKeys, key)
				seen[key] = true
			}
		}
		return uniqueKeys, nil
	}

	// If no results from either cache, return the primary error if available
	if primaryErr != nil {
		return nil, primaryErr
	}
	return allKeys, fallbackErr
}

// GetJSON retrieves a value from cache and unmarshals it to the provided struct
func (ac *AdaptiveCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	ac.mutex.RLock()
	usingFallback := ac.usingFallback
	ac.mutex.RUnlock()

	if usingFallback || ac.redisCache == nil {
		return ac.fallbackCache.GetJSON(ctx, key, dest)
	}

	err := ac.redisCache.GetJSON(ctx, key, dest)
	if err != nil {
		// If Redis fails, try fallback cache
		log.Printf("Redis error during GetJSON, trying fallback: %v", err)
		ac.mutex.Lock()
		ac.usingFallback = true
		ac.mutex.Unlock()

		return ac.fallbackCache.GetJSON(ctx, key, dest)
	}

	return err
}

// SetJSON sets a value in cache with expiration after marshaling to JSON
func (ac *AdaptiveCache) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ac.mutex.RLock()
	usingFallback := ac.usingFallback
	ac.mutex.RUnlock()

	var primaryErr error
	var fallbackErr error

	// Try primary cache first if it's available and we're not using fallback
	if !usingFallback && ac.redisCache != nil {
		primaryErr = ac.redisCache.SetJSON(ctx, key, value, expiration)
		if primaryErr != nil {
			// Switch to fallback on error
			log.Printf("Redis error during SetJSON, switching to fallback: %v", primaryErr)
			ac.mutex.Lock()
			ac.usingFallback = true
			ac.mutex.Unlock()
		} else {
			// Set in fallback cache as well for consistency when Redis is working
			ac.fallbackCache.SetJSON(ctx, key, value, expiration)
			return nil // Success with primary cache
		}
	}

	// Set in fallback cache as well for consistency
	fallbackErr = ac.fallbackCache.SetJSON(ctx, key, value, expiration)
	if fallbackErr != nil {
		log.Printf("Fallback cache error during SetJSON: %v", fallbackErr)
	}

	// Return the primary error if we were trying to use Redis and it failed
	if primaryErr != nil {
		return primaryErr
	}

	return fallbackErr
}

// Close closes the fallback cache
func (ac *AdaptiveCache) Close() {
	ac.fallbackCache.Close()
}