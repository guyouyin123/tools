package qredis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// BatchHGetAll 获取在哈希表中指定 keys 的所有字段和值。
func (p *Redis) BatchHGetAll(ctx context.Context, keys ...string) (res map[string]map[string]string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	cmds := make(map[string]*redis.StringStringMapCmd, len(keys))
	if _, err := p.client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			cmds[key] = pipe.HGetAll(ctx, key)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	result := make(map[string]map[string]string, len(cmds))
	for _, key := range keys {
		if v, ok := cmds[key]; ok && v != nil {
			result[key] = v.Val()
		}
	}
	return result, nil
}

// HMGet 获取所有给定字段的值
func (p *Redis) HMGet(ctx context.Context, key string, fields ...string) (res []interface{}, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	vals := make([]interface{}, len(fields))
	if vals, err = p.client.HMGet(ctx, key, fields...).Result(); err != nil {
		return res, err
	}
	return vals, err
}

// HMSet 同时将多个 field-value (域-值)对设置到哈希表 key 中。
func (p *Redis) HMSet(ctx context.Context, key string, fields map[string]interface{}) (res bool, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.HMSet(ctx, key, fields).Result()
}

// HSetNX 只有在字段 field 不存在时，设置哈希表字段的值。
func (p *Redis) HSetNX(ctx context.Context, key, field string, value interface{}) (res bool, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.HSetNX(ctx, key, field, value).Result()
}

// HGet 获取存储在哈希表中指定字段的值。
func (p *Redis) HGet(ctx context.Context, key, field string) (res string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.HGet(ctx, key, field).Result()
}

// HGetAll 获取在哈希表中指定 key 的所有字段和值。
func (p *Redis) HGetAll(ctx context.Context, key string) (res map[string]string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.HGetAll(ctx, key).Result()
}

// HSet 将哈希表 key 中的字段 field 的值设为 value。
func (p *Redis) HSet(ctx context.Context, key, field string, value interface{}) (res int64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.HSet(ctx, key, field, value).Result()
}

// HDel 删除一个或多个哈希表字段
func (p *Redis) HDel(ctx context.Context, key, field string) (res int64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.HDel(ctx, key, field).Result()
}

// HIncrBy 为哈希表 key 中的指定字段的整数值加上增量 increment 。
func (p *Redis) HIncrBy(ctx context.Context, key, field string, incr int64) (res int64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.HIncrBy(ctx, key, field, incr).Result()
}

// HLen 获取哈希表中字段的数量
func (p *Redis) HLen(ctx context.Context, key string) (res int64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.HLen(ctx, key).Result()
}

// HKeys 获取所有哈希表中的字段
func (p *Redis) HKeys(ctx context.Context, key string) (res []string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.HKeys(ctx, key).Result()
}

// HScan 迭代哈希表中的键值对
func (p *Redis) HScan(ctx context.Context, key string, cursorR uint64, match string, count int64) (keys []string, cursor uint64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		err = errRedis
		return
	}
	return p.client.HScan(ctx, key, cursorR, match, count).Result()
}
