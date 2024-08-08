package qredis

import (
	"context"
)

// EvalSha 缓存执行lua脚本
func (p *Redis) EvalSha(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	// 安全检查
	if p == nil || p.client == nil {
		return nil, errRedis
	}
	// 将 Lua 脚本缓存到 Redis 服务器
	scriptSha, err := p.client.ScriptLoad(ctx, script).Result()
	if err != nil {
		return nil, err
	}

	return p.client.EvalSha(ctx, scriptSha, keys, args...).Result()
}

// Eval 直接执行lua脚本
func (p *Redis) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	// 安全检查
	if p == nil || p.client == nil {
		return nil, errRedis
	}
	return p.client.Eval(ctx, script, keys, args...).Result()
}
