package cache

import "github.com/go-redis/redis/v8"

// init redisclient
func NewRedisClient(addr, pw string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
}
