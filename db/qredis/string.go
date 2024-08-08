package qredis

import (
	"context"
	"time"
)

// Append 如果 key 已经存在并且是一个字符串,Append 将指定的 value 追加到该 key 原来值（value）的末尾。
func (p *Redis) Append(ctx context.Context, key, value string) (res int64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.Append(ctx, key, value).Result()
}

// Get 获取指定 key 的值。
func (p *Redis) Get(ctx context.Context, key string) (res string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.Get(ctx, key).Result()
}

// MGet 获取所有(一个或多个)给定 key 的值。
func (p *Redis) MGet(ctx context.Context, key ...string) (res []interface{}, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.MGet(ctx, key...).Result()
}

// Set 设置指定 key 的值,expiration=0就是永久
func (p *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (res string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.Set(ctx, key, value, expiration).Result()
}

// MSet 同时设置一个或多个 key-value 对。
// MSet(map[string]interface{}{"key1": "value1", "key2": "value2"})
func (p *Redis) MSet(ctx context.Context, values ...interface{}) (res string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.MSet(ctx, values).Result()
}

// IncrBy 将 key 所储存的值加上给定的增量值（increment）。
func (p *Redis) IncrBy(ctx context.Context, key string, value int64) (res int64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.IncrBy(ctx, key, value).Result()
}
