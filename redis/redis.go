package redis

import (
	"github.com/go-redis/redis"
)

type RedisConfigType struct {
	Addr     string
	Password string
	DB       int
}

func GetRedisClient(redisConfig *RedisConfigType) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	return client
}
