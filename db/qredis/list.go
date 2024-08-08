package qredis

import (
	"context"
	"time"
)

// LPop 移出并获取列表的第一个元素
func (p *Redis) LPop(ctx context.Context, key string) (res string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.LPop(ctx, key).Result()
}

// RPop 移出并获取列表的最后一个元素
func (p *Redis) RPop(ctx context.Context, key string) (res string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.RPop(ctx, key).Result()
}

// BLPop 移出并获取列表的第一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止
func (p *Redis) BLPop(ctx context.Context, timeout time.Duration, keys ...string) (res []string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.BLPop(ctx, timeout, keys...).Result()
}

// BRPop 移出并获取列表的最后一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
func (p *Redis) BRPop(ctx context.Context, timeout time.Duration, keys ...string) (res []string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.BRPop(ctx, timeout, keys...).Result()
}

// LPush 将一个或多个值插入到列表头部
func (p *Redis) LPush(ctx context.Context, key string, values ...interface{}) (res int64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.LPush(ctx, key, values).Result()
}

// RPush 将一个或多个值插入到列表尾部
func (p *Redis) RPush(ctx context.Context, key string, values ...interface{}) (res int64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.RPush(ctx, key, values).Result()
}

// LRange 获取列表指定下标索引范围内的元素
// eg:消费者一次取多个
func (p *Redis) LRange(ctx context.Context, key string, start, stop int64) (res []string, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.LRange(ctx, key, start, stop).Result()
}

// LLen 获取列表长度
func (p *Redis) LLen(ctx context.Context, key string) (res int64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.LLen(ctx, key).Result()
}
