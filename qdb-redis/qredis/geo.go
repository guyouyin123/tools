package qredis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

//Redis GEO 主要用于存储地理位置信息，并对存储的信息进行操作，该功能在 Redis 3.2 版本新增。

// GeoAdd 添加地理位置的坐标。
func (p *Redis) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) (res int64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.GeoAdd(ctx, key, geoLocation...).Result()
}

// GeoPos 获取地理位置的坐标。
func (p *Redis) GeoPos(ctx context.Context, key string, members ...string) (res []*redis.GeoPos, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.GeoPos(ctx, key, members...).Result()
}

// GeoDist 计算两个位置之间的距离。
// member1 member2 为两个地理位置。
// unit m:米，默认单位。km:千米。mi:英里。ft:英尺。
func (p *Redis) GeoDist(ctx context.Context, key string, member1, member2, unit string) (res float64, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.GeoDist(ctx, key, member1, member2, unit).Result()
}

// GeoRadius 根据用户给定的经纬度坐标来获取指定范围内的地理位置集合。
func (p *Redis) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) (res []redis.GeoLocation, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.GeoRadius(ctx, key, longitude, latitude, query).Result()
}

// GeoRadiusByMember 根据储存在位置集合里面的某个地点获取指定范围内的地理位置集合。
func (p *Redis) GeoRadiusByMember(ctx context.Context, key, member string, query *redis.GeoRadiusQuery) (res []redis.GeoLocation, err error) {
	// 安全检查
	if p == nil || p.client == nil {
		return res, errRedis
	}
	return p.client.GeoRadiusByMember(ctx, key, member, query).Result()
}
