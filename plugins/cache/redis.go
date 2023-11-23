package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type redisClient struct {
	Client  *redis.Client
	Context context.Context
}

var (
	once  sync.Once
	Redis *redisClient
)

// @author: lipper
// @object: *redis.Client
// @function: NewRedis
// @description: 实例redis
// @return: *redis.Client
func NewRedis(conf RedisConf) *redis.Client {
	once.Do(func() {
		Redis = connectRedis(conf)
	})
	return Redis.Client
}

// 创建新的redis连接
func connectRedis(conf RedisConf) *redisClient {
	if err := conf.Validate(); err != nil {
		panic(err)
	}
	rds := &redisClient{}
	rds.Context = context.Background()
	// 初始化NewClient 连接
	rds.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})
	// 测试链接
	if err := rds.Ping(); err != nil {
		panic(err)
	}
	return rds
}

// Ping 用以测试 redis 连接是否正常
func (rds redisClient) Ping() error {
	_, err := rds.Client.Ping(rds.Context).Result()
	return err
}

// Set 存储 key 对应的 value，且设置 expiration 过期时间
func (rds redisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return rds.Client.Set(rds.Context, key, value, expiration).Err()
}

// Get 获取 key 对应的 value
func (rds redisClient) Get(key string) (string, error) {
	result, err := rds.Client.Get(rds.Context, key).Result()
	return result, err
}

// Has 判断一个 key 是否存在，内部错误和 redis.Nil 都返回 false
func (rds redisClient) Has(key string) (string, error) {
	result, err := rds.Client.Get(rds.Context, key).Result()
	return result, err
}

// Del 删除存储在 redis 里的数据，支持多个 key 传参
func (rds redisClient) Del(keys ...string) error {
	err := rds.Client.Del(rds.Context, keys...).Err()
	return err
}

// FlushDB 清空当前 redis db 里的所有数据
func (rds redisClient) FlushDB() error {
	err := rds.Client.FlushDB(rds.Context).Err()
	return err
}
