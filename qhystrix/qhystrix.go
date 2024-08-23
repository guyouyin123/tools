package qhystrix

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"time"
)

func InitConf(Timeout, MaxConcurrentRequests, SleepWindow, ErrorPercentThreshold, RequestVolumeThreshold int) *hystrix.CommandConfig {
	config := &hystrix.CommandConfig{
		Timeout:                Timeout,                //执行command的超时时间
		MaxConcurrentRequests:  MaxConcurrentRequests,  //command的最大并发量
		SleepWindow:            SleepWindow,            //过多长时间，熔断器再次检测是否开启。单位毫秒
		ErrorPercentThreshold:  ErrorPercentThreshold,  //错误率 请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动
		RequestVolumeThreshold: RequestVolumeThreshold, //请求阈值(一个统计窗口10秒内请求数量) 熔断器是否打开首先要满足这个条件；这里的设置表示至少有5个请求才进行ErrorPercentThreshold错误百分比计算
	}
	return config
}

func Hystrix(conf *hystrix.CommandConfig) {
	hystrix.ConfigureCommand("test", *conf)
	cbs, _, _ := hystrix.GetCircuit("test")
	defer hystrix.Flush()
	i := 1
	for {
		start1 := time.Now()
		hystrix.Do("test", RunFuncTest, func(e error) error {
			fmt.Println("服务器错误 触发 fallbackFunc 调用函数执行失败 : ", e)
			return nil
		})
		fmt.Println("请求次数:", i+1, ";用时:", time.Now().Sub(start1), ";熔断器开启状态:", cbs.IsOpen(), "请求是否允许：", cbs.AllowRequest())
		time.Sleep(time.Second)
		i++
	}
}

// 业务函数
func RunFuncTest() error {
	time.Sleep(time.Second * 3)
	return nil
}
