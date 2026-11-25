package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"ticketBooking/internal/config"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client  *redis.Client
	enabled bool
	mu      sync.RWMutex
}

var (
	defaultCache *Cache
	once         sync.Once
)

// Init initializes the singleton Redis cache instance
func Init(cfg config.Config) *Cache {
	once.Do(func() {
		defaultCache = New(cfg)
	})
	return defaultCache
}

// Get returns the default cache instance (falls back to a disabled no-op cache if not initialized)
func Get() *Cache {
	if defaultCache == nil {
		return &Cache{enabled: false}
	}
	return defaultCache
}

// New creates a new Cache instance with graceful degradation
func New(cfg config.Config) *Cache {
	redisURL := strings.TrimSpace(cfg.RedisURL)
	redisAddr := strings.TrimSpace(cfg.RedisAddr)

	var opt *redis.Options

	if redisURL != "" {
		parsedOpt, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Printf("[Redis] Failed to parse REDIS_URL (%s): %v. Cache operating in no-op mode.\n", redisURL, err)
			return &Cache{enabled: false}
		}
		opt = parsedOpt
	} else if redisAddr != "" {
		opt = &redis.Options{
			Addr: redisAddr,
		}
	} else {
		// Default to localhost:6379 in development if neither is provided
		opt = &redis.Options{
			Addr: "localhost:6379",
		}
	}

	opt.DialTimeout = 3 * time.Second
	opt.ReadTimeout = 2 * time.Second
	opt.WriteTimeout = 2 * time.Second

	rdb := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("[Redis] Could not connect to Redis at %s (%v). Operating in cache-bypass mode.\n", opt.Addr, err)
		return &Cache{client: rdb, enabled: false}
	}

	log.Printf("[Redis] Connected successfully to Redis at %s\n", opt.Addr)
	return &Cache{
		client:  rdb,
		enabled: true,
	}
}

// IsEnabled returns whether Redis cache is currently active
func (c *Cache) IsEnabled() bool {
	if c == nil {
		return false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.enabled
}

// Get retrieves an item from cache and unmarshals it into dest.
// Returns (true, nil) on Cache Hit, (false, nil) on Cache Miss or error.
func (c *Cache) Get(ctx context.Context, key string, dest any) (bool, error) {
	if !c.IsEnabled() {
		return false, nil
	}

	val, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil // Cache miss
	}
	if err != nil {
		log.Printf("[Redis] Get error for key '%s': %v\n", key, err)
		return false, nil // Treat network/redis error as cache miss to avoid breaking the caller
	}

	if err := json.Unmarshal(val, dest); err != nil {
		log.Printf("[Redis] Unmarshal error for key '%s': %v\n", key, err)
		return false, nil
	}

	return true, nil
}

// Set serializes value to JSON and stores it with the given TTL
func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if !c.IsEnabled() {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("[Redis] Marshal error for key '%s': %v\n", key, err)
		return err
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

// Delete removes specific key(s)
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	if !c.IsEnabled() || len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

// DeleteByPrefix safely finds and deletes all keys matching `prefix*` using SCAN (non-blocking)
func (c *Cache) DeleteByPrefix(ctx context.Context, prefix string) error {
	if !c.IsEnabled() {
		return nil
	}

	var cursor uint64
	matchPattern := prefix + "*"

	for {
		keys, nextCursor, err := c.client.Scan(ctx, cursor, matchPattern, 100).Result()
		if err != nil {
			log.Printf("[Redis] Scan error for prefix '%s': %v\n", prefix, err)
			return err
		}

		if len(keys) > 0 {
			if delErr := c.client.Del(ctx, keys...).Err(); delErr != nil {
				log.Printf("[Redis] Del error during prefix invalidation '%s': %v\n", prefix, delErr)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}

// GenerateKey produces a readable, deterministic cache key from a prefix and query components
func GenerateKey(prefix string, parts ...any) string {
	var sb strings.Builder
	sb.WriteString(prefix)
	for _, p := range parts {
		str := fmt.Sprint(p)
		if str != "" {
			sb.WriteString(":")
			sb.WriteString(str)
		}
	}
	key := sb.String()
	if len(key) > 150 {
		hash := sha256.Sum256([]byte(key))
		return fmt.Sprintf("%s:%s", prefix, hex.EncodeToString(hash[:])[:16])
	}
	return key
}
