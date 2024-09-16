package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	Rdb *redis.Client
}

var (
	instanceRedis *RedisClient
	onceRedis     sync.Once
)

func GetRedisClient() *RedisClient {
	return instanceRedis
}

func InitRedisClient() {
	onceRedis.Do(func() {
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

		instanceRedis = &RedisClient{
			Rdb: Client,
		}
	})
}

func (rdc *RedisClient) Set(key string, value string) error {
	err := rdc.SetWithExpired(key, value, 0)
	if err != nil {
		return err
	}

	return nil
}

func (rdc *RedisClient) SetWithExpired(key string, value string, expired time.Duration) error {
	err := rdc.Rdb.Set(ctx, key, value, expired).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rdc *RedisClient) Get(key string) (string, error) {
	result, err := rdc.Rdb.Get(ctx, key).Result()

	switch {
	case err == redis.Nil:
		fmt.Println("Keys does not exists")
		return "", redis.Nil
	case err != nil:
		fmt.Println("Failed to get data")
		return "", err
	case result == "":
		fmt.Println("Value empty")
		return "", fmt.Errorf("Value empty")
	}

	return result, nil
}
