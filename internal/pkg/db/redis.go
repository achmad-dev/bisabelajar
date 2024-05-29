package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(addr string, password string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	// Ping the Redis server to check the connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
