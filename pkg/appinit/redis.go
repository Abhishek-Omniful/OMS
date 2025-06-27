package appinit

import (
	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/redis"
)

var RedisClient *redis.Client

func ConnectRedis() *redis.Client {
	ctx := mycontext.GetContext()

	logger.Println(i18n.Translate(ctx, "Connecting to Redis..."))

	config := &redis.Config{
		Hosts:       []string{"127.0.1.1:6379"},
		PoolSize:    50,
		MinIdleConn: 10,
		DB:          0,
	}

	RedisClient = redis.NewClient(config)

	logger.Infof(i18n.Translate(ctx, "Redis initialized successfully!"))

	return RedisClient
}

func GetRedis() *redis.Client {
	return RedisClient
}
