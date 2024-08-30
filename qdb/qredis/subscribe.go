package qredis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// PSubscribe 订阅过期消息
// 需要在 Redis 配置文件中开启过期消息通知功能--conf中添加配置 notify-keyspace-events Ex
func (p *Redis) PSubscribe(ctx context.Context, do func(message *redis.Message), patterns ...string) error {
	// 安全检查
	if p == nil || p.client == nil {
		return errRedis
	}
	pubsub := p.client.PSubscribe(ctx, patterns...)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return err
	}
	ch := pubsub.Channel()
	for msg := range ch {
		do(msg)
	}
	return nil
}
