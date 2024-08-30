package qredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var (
	errRedis        = fmt.Errorf("redis is nil")
	errRedisClient  = fmt.Errorf("redis client is nil")
	errRedisCClient = fmt.Errorf("redis cclient is nil")
	Nil             = redis.Nil
)

type Redis struct {
	client *redis.Client
	//cclient *redis.ClusterClient //集群模式占不开启
}

func (p *Redis) InitRedis(conf *redis.Options) error {
	//conf = &redis.Options{
	//	Addr: conf.Addr,
	//	//Username: conf.UserName,
	//	Password:     conf.Password,
	//	DB:           2,
	//	MaxRetries:   3,
	//	MinIdleConns: 8,
	//	PoolTimeout:  2 * time.Minute,
	//	IdleTimeout:  10 * time.Minute,
	//	ReadTimeout:  2 * time.Minute,
	//	WriteTimeout: 1 * time.Minute,
	//}
	redisCli := redis.NewClient(conf)
	_, err := redisCli.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	p.client = redisCli
	return nil
}
