package appinit

import (
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/redis"
)

var RedisClient *redis.Client

func ConnectRedis() *redis.Client {
	log.Println("Connecting to Redis...")
	config := &redis.Config{
		Hosts:       []string{"127.0.1.1:6379"},
		PoolSize:    50,
		MinIdleConn: 10,
		DB:          0,
	}
	RedisClient = redis.NewClient(config)
	log.Infof("Redis initialized successfully!")
	return RedisClient
}

func GetRedis() *redis.Client {
	return RedisClient
}
