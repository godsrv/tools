package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	once   sync.Once
	Client *redis.Client
)

// @author: lipper
// @object: *redis.Client
// @function: NewRedis
// @description: 实例redis
// @return: *redis.Client
func NewRedis(conf RedisConf) *redis.Client {
	once.Do(func() {
		Client = connectRedis(conf)
	})
	return Client
}

// 创建新的redis连接
func connectRedis(conf RedisConf) *redis.Client {
	if err := conf.Validate(); err != nil {
		panic(err)
	}
	ct := context.Background()
	// 初始化NewClient 连接
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})
	// 测试链接
	_, err := Client.Ping(ct).Result()
	if err != nil {
		logrus.Panicf("ping redis err: %v", err)
	}
	logrus.Info("connect to redis success")
	return Client
}
