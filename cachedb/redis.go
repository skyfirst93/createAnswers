package cachedb

import (
	"log"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func InitRedis() bool {
	redisAddress := "127.0.0.1:6379"
	redisPassword := "password01"
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       0,
		PoolSize: 4000,
	})
	_, err := RedisClient.Ping().Result()

	if err != nil {
		log.Printf("InitRedis: failed to connect to redis %v", err)
		return false
	}

	return true
}
