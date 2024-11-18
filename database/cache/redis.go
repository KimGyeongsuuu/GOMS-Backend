package cache

import (
	"context"
	"fmt"
	"log"

	"GOMS-BACKEND-GO/global/config"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(ctx context.Context) *redis.Client {
	host := config.Data().Redis.Host
	port := config.Data().Redis.Port
	addr := fmt.Sprintf("%s:%d", host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
		return nil
	}
	return rdb
}
