package db

import (
	"context"
	"fmt"
	"log"
	"sync"

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
	err := rdc.Rdb.Set(ctx, key, value, 0).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rdc *RedisClient) Get(key string) (string, error) {
	result, err := rdc.Rdb.Get(ctx, key).Result()

	if err != nil {
		return result, err
	}

	return result, nil
}
