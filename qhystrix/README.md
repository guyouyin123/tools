# 熔断工具



```go

import (
	"testing"
)
func TestHystrix(t *testing.T) {
	//config := hystrix.CommandConfig{
	//	Timeout:                2000, //执行command的超时时间
	//	MaxConcurrentRequests:  8,    //command的最大并发量
	//	SleepWindow:            2000, //过多长时间，熔断器再次检测是否开启。单位毫秒
	//	ErrorPercentThreshold:  30,   //错误率 请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动
	//	RequestVolumeThreshold: 5,    //请求阈值(一个统计窗口10秒内请求数量) 熔断器是否打开首先要满足这个条件；这里的设置表示至少有5个请求才进行ErrorPercentThreshold错误百分比计算
	//}
	config := InitConf(2000, 8, 2000, 30, 5)
	Hystrix(config)
}

```