package qredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

var p Redis

func InitDemo() {
	conf := &redis.Options{
		Addr: "127.0.0.1:6379",
		//Username: conf.UserName,
		//Password:     conf.Password,
		DB:           0,
		MaxRetries:   3,
		MinIdleConns: 8,
		PoolTimeout:  2 * time.Minute,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	}
	err := p.InitRedis(conf)
	if err != nil {
		panic(err)
	}
}

func TestRedis_hash(t *testing.T) {
	InitDemo()
	ctx := context.Background()
	//boo, _ := p.HMSet(ctx, "user2", map[string]interface{}{"name": "Jeff2", "age2": 18})
	//fmt.Println(boo)
	//res, _ := p.BatchHGetAll(ctx, "user1", "user2")
	//fmt.Println(res)

	p.Set(ctx, "key", "value", 0)
	res, _ := p.Get(ctx, "key")
	fmt.Println(res)
}

func TestRedis_geo(t *testing.T) {
	InitDemo()
	ctx := context.Background()
	geoList := make([]*redis.GeoLocation, 0)
	shanghai := &redis.GeoLocation{
		Name:      "shanghai",
		Longitude: 121.472644, //经度
		Latitude:  30.404,     //纬度
		Dist:      100,
		GeoHash:   100,
	}
	beijing := &redis.GeoLocation{
		Name:      "beijing",
		Longitude: 116.23,
		Latitude:  39.54,
		Dist:      100,
		GeoHash:   100,
	}
	geoList = append(geoList, shanghai)
	geoList = append(geoList, beijing)
	//p.GeoAdd(ctx, "position", geoList...)
	//res, _ := p.GeoPos(ctx, "position", "shanghai", "beijing")
	//fmt.Println(res)
	//res, _ := p.GeoDist(ctx, "position", "shanghai", "beijing", "m")
	//res2, _ := p.GeoDist(ctx, "position", "shanghai", "beijing", "km")
	//fmt.Println(res)
	//fmt.Println(res2)

	query := &redis.GeoRadiusQuery{
		Radius: 10000, // 1km
		Unit:   "km",
	}
	//res2, _ := p.GeoRadius(ctx, "position", 121.472644, 30.404, query)
	//fmt.Println(res2)

	res3, _ := p.GeoRadiusByMember(ctx, "position", "shanghai", query) //shanghai范围内10000km内的点
	fmt.Println(res3)
}

func TestRedis_List(t *testing.T) {
	InitDemo()
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		_, err := p.LPush(ctx, "list_test", i)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	res, _ := p.LRange(ctx, "list_test", 0, 10)
	fmt.Println(res)
}

func TestRedis_Lua(t *testing.T) {
	InitDemo()
	ctx := context.Background()
	script := `
		local key = KEYS[1]
		local value = ARGV[1]
		redis.call('SET', key, value)
		return 'OK'
	`
	//res, _ := p.EvalSha(ctx, script, []string{"lock3"}, []string{"value"})
	res, _ := p.Eval(ctx, script, []string{"lock4"}, []string{"value4"})
	fmt.Println(res)

}
