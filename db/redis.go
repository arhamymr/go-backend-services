package db

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Rdb *redis.Client
}

func NewRedisClient() *RedisClient {
	Client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Host and port
		Password: "",               // Password
		DB:       0,                // Default DB; this part is not specified in the URL and typically defaults to 0
	})

	pong, err := Client.Ping(context.Background()).Result()

	if err != nil {
		log.Fatalf("Could not connect to redis %v \n", err)
	}

	fmt.Println(pong, "Successfully connected to Redis")
	return &RedisClient{
		Rdb: Client,
	}
}
