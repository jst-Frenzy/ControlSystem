package dataBase

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

var RedisDB *redis.Client

func InitRedis() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ping, err := RedisDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis")
	}
	fmt.Println(ping)
}
