# 过期安全map


过期map
特点：
1.设置key时必须设置key的存在时间
2.加锁的安全map

优化：过期key的清理
1.在获取key时，判断距离上次清理时间超过60秒则启动协程清理一次
2.无需启动一个永久协程循环检查清理失效的key，浪费资源


```go

import (
"fmt"
"time"
qexmap "github.com/guyouyin123/tools/qexmap"
)

func TestNewExpiringMap(t *testing.T) {
	type user struct {
		IdCard string
		Name   string
	}
	jeff := &user{
		IdCard: "123456",
		Name:   "jeff",
	}

	userIdMap := qexmap.NewExpiringMap()

	userIdMap.Set(jeff.IdCard, jeff, time.Second*10)
	v, ok := userIdMap.Get(jeff.IdCard)
	if ok {
		fmt.Println(v)
	}
}


```

