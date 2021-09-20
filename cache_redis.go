package d

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

var (
	Cache cache
	CacheClient *redis.Client
)

type cache struct {
	Addr string
	Password string
	DB int
	Optional optionalRedis
	enable bool
}

type optionalRedis struct {

}

// https://github.com/go-redis/redis
func (this *cache) Init() error {
	CacheClient = redis.NewClient(&redis.Options{
		Addr:     this.Addr,
		Password: this.Password, // no password set
		DB:       this.DB,  // use default DatabaseClient
	})

	_, err := CacheClient.Ping(context.TODO()).Result()
	if err != nil {
		return errors.New("redis connection failed")
	}

	this.enable = true
	return nil
}

// 获取启动的状态
func (this *cache) GetEnabledStatus() bool {
	return this.enable
}

func (this *cache) CacheClearAll() {
	if this.enable {
		CacheClient.FlushDB(context.Background())
	}
	return
}