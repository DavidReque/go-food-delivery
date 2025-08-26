package redis

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

const (
	maxRetries      = 5                      // retry 5 times
	minRetryBackoff = 300 * time.Millisecond // retry backoff 300ms
	maxRetryBackoff = 500 * time.Millisecond
	dialTimeout     = 5 * time.Second // dial timeout 5s
	readTimeout     = 5 * time.Second // read timeout 5s
	writeTimeout    = 3 * time.Second // write timeout 3s
	minIdleConns    = 20              // min idle connections 20
	poolTimeout     = 6 * time.Second // pool timeout 6s
)

func NewRedisClient(cfg *RedisOptions) *redis.Client {
	universalClient := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username:        cfg.Username, // Redis Cloud username
		Password:        cfg.Password, // Redis Cloud password
		DB:              cfg.Database, // use default database
		MaxRetries:      maxRetries,
		MinRetryBackoff: minRetryBackoff,
		MaxRetryBackoff: maxRetryBackoff,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    minIdleConns,
		PoolTimeout:     poolTimeout,
	})

	if cfg.EnableTracing {
		_ = redisotel.InstrumentTracing(universalClient)
	}

	return universalClient
}
