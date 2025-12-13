package cache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/AndikaPrasetia/pos-cafee/internal/config"
)

// RedisManager manages Redis connections with automatic reconnection capability
type RedisManager struct {
	client    *redis.Client
	config    *config.AppConfig
	mu        sync.RWMutex
	isHealthy bool
}

// NewRedisManager creates a new Redis manager with automatic reconnection capabilities
func NewRedisManager(cfg *config.AppConfig) (*RedisManager, error) {
	rm := &RedisManager{
		config:    cfg,
		isHealthy: false,
	}
	
	client, err := rm.createClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %v", err)
	}
	
	rm.client = client
	rm.isHealthy = true
	
	return rm, nil
}

// createClient creates a new Redis client based on configuration
func (rm *RedisManager) createClient(cfg *config.AppConfig) (*redis.Client, error) {
	if cfg.Redis.URL == "" {
		return nil, fmt.Errorf("REDIS_URL is not set in configuration")
	}

	opts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	// Configure connection pooling and timeouts
	opts.PoolSize = 20                 // Connection pool size
	opts.MinIdleConns = 5              // Minimum idle connections
	opts.PoolTimeout = 30 * time.Second // Connection pool timeout
	opts.IdleTimeout = 300 * time.Second // Connection idle timeout
	opts.IdleCheckFrequency = 60 * time.Second // Frequency of idle connection checks
	opts.DialTimeout = 10 * time.Second  // Dial timeout
	opts.ReadTimeout = 10 * time.Second  // Read timeout
	opts.WriteTimeout = 10 * time.Second // Write timeout

	client := redis.NewClient(opts)
	
	// Test the connection with retry logic
	if err := rm.testConnection(client); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to Redis")
	return client, nil
}

// testConnection tests if Redis connection is working
func (rm *RedisManager) testConnection(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}

// GetClient returns the Redis client, automatically reconnecting if needed
func (rm *RedisManager) GetClient() *redis.Client {
	rm.mu.RLock()
	client := rm.client
	isHealthy := rm.isHealthy
	rm.mu.RUnlock()

	if !isHealthy {
		rm.mu.Lock()
		// Check again after acquiring write lock
		if !rm.isHealthy {
			log.Println("Attempting to reconnect to Redis...")
			
			// Attempt to create a new client
			newClient, err := rm.createClient(rm.config)
			if err != nil {
				log.Printf("Failed to reconnect to Redis: %v", err)
				rm.mu.Unlock()
				return nil
			}

			// Close the old client if it exists
			if rm.client != nil {
				rm.client.Close()
			}

			rm.client = newClient
			rm.isHealthy = true
			log.Println("Successfully reconnected to Redis")
		}
		client = rm.client
		rm.mu.Unlock()
	}

	return client
}

// IsHealthy returns whether the Redis connection is healthy
func (rm *RedisManager) IsHealthy() bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.isHealthy
}

// Monitor continuously monitors Redis connection health
func (rm *RedisManager) Monitor(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Redis monitor...")
			return
		case <-ticker.C:
			rm.checkHealth()
		}
	}
}

// checkHealth checks the health of Redis connection
func (rm *RedisManager) checkHealth() {
	rm.mu.RLock()
	client := rm.client
	rm.mu.RUnlock()

	if client == nil {
		rm.setHealth(false)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Redis health check failed: %v", err)
		rm.setHealth(false)
	} else {
		rm.setHealth(true)
	}
}

// setHealth sets the health status of Redis connection
func (rm *RedisManager) setHealth(healthy bool) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.isHealthy = healthy
}

// GetMetrics gets Redis metrics for monitoring
func (rm *RedisManager) GetMetrics() (map[string]interface{}, error) {
	client := rm.GetClient()
	if client == nil {
		return nil, fmt.Errorf("Redis client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get Redis info
	info, err := client.Info(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis info: %v", err)
	}

	// Parse info to extract metrics
	metrics := make(map[string]interface{})
	
	// Basic metrics
	lines := splitLines(info)
	for _, line := range lines {
		if keyValue := splitKeyValue(line); keyValue != nil {
			metrics[keyValue[0]] = keyValue[1]
		}
	}

	// Additional computed metrics
	connectionCount, ok := metrics["connected_clients"].(string)
	if ok {
		metrics["connected_clients_count"], _ = fmt.Sscanf(connectionCount, "%d", &connectionCount)
	}

	totalCommandsProcessed, ok := metrics["total_commands_processed"].(string)
	if ok {
		metrics["total_commands_processed_count"], _ = fmt.Sscanf(totalCommandsProcessed, "%d", &totalCommandsProcessed)
	}

	usedMemory, ok := metrics["used_memory"].(string)
	if ok {
		metrics["used_memory_bytes"], _ = fmt.Sscanf(usedMemory, "%d", &usedMemory)
	}

	return metrics, nil
}

// Helper function to split string by lines
func splitLines(s string) []string {
	var lines []string
	currentLine := ""
	for _, char := range s {
		if char == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else {
			currentLine += string(char)
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

// Helper function to split key-value pairs by ':'
func splitKeyValue(s string) []string {
	parts := []string{}
	currentPart := ""
	inQuote := false

	for _, char := range s {
		if char == '"' {
			inQuote = !inQuote
		} else if char == ':' && !inQuote {
			parts = append(parts, currentPart)
			currentPart = ""
		} else {
			currentPart += string(char)
		}
	}

	if currentPart != "" {
		parts = append(parts, currentPart)
	}

	if len(parts) != 2 {
		return nil
	}

	return parts
}