package qredis

import (
	"context"
	"time"
)

// Keys 查看当前所有有效key
func (p *Redis) Keys(ctx context.Context, key string) (res []string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.Keys(ctx, key).Result()
}

// Expire 设置key过期时间
func (p *Redis) Expire(ctx context.Context, key string, expiration time.Duration) (res bool, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.Expire(ctx, key, expiration).Result()
}

// Exists 检查给定 key 是否存在。存在返回1，否则返回0
func (p *Redis) Exists(ctx context.Context, keys ...string) (res int64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.Exists(ctx, keys...).Result()
}

// Del 该命令用于在 key 存在时删除 key。
func (p *Redis) Del(ctx context.Context, keys ...string) (res int64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.Del(ctx, keys...).Result()
}

// Persist 移除 key 的过期时间，key 将持久保持。
func (p *Redis) Persist(ctx context.Context, key string) (res bool, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.Persist(ctx, key).Result()
}

// TTL 以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)
func (p *Redis) TTL(ctx context.Context, key string) (res int, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	ttl, err := p.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(ttl.Seconds()), nil
}

// RenameNX 仅当 newKey 不存在时，将 key 改名为 newKey 。
func (p *Redis) RenameNX(ctx context.Context, key, newKey string) (res bool, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.RenameNX(ctx, key, newKey).Result()
}
