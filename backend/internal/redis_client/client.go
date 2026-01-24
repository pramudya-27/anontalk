package redis_client

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() *RedisClient {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		fmt.Printf("Warning: Could not connect to Redis: %v\n", err)
	} else {
		fmt.Println("Connected to Redis")
	}

	return &RedisClient{Client: rdb}
}
