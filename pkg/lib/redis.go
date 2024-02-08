package lib

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func InitializeAndTestRedisClient() (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr, should be set via config / flags
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	status := redisClient.Ping(context.Background())
	if status.Err() != nil {
		return nil, status.Err()
	}
	return redisClient, nil
}
