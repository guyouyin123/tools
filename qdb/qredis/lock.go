package qredis

import (
	"context"
	"time"
)

// Lock 锁
func (p *Redis) Lock(ctx context.Context, key string, expire int) (bool, error) {
	now := time.Now().Unix()
	boo, err := p.client.SetNX(ctx, key, now, time.Duration(expire)*time.Second).Result()
	if err != nil {
		return false, err
	}
	return boo, nil
}

// UnLock 释放锁
func (p *Redis) UnLock(ctx context.Context, key string) (bool, error) {
	_, err := p.client.Del(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

// LockTTL 获取锁的剩余时间--以秒为单位
func (p *Redis) LockTTL(ctx context.Context, key string) (int, error) {
	ttl, err := p.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(ttl.Seconds()), nil
}
