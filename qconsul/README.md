# qconsul-Consul--本文实现客户端

是 HashiCorp 开发的一款工具，主要用于服务发现、配置管理和分布式系统的健康检查。它允许微服务之间进行通信，并提供一种简便的方法来管理服务的配置和状态。

**服务端官方下载链接**：https://developer.hashicorp.com/consul/install

后台启动命令：nohup ./consul agent -dev &



```go
import (
	"fmt"
	"testing"
)

func TestRunTest(t *testing.T) {
	Address := "127.0.0.1:8500"
	//初始化
	consulClient, err := InitConn(Address)
	if err != nil {
		return
	}

	type HReturnFeeRefreshReq struct {
		IntvDtSta string //面试日期开始时间
		IntvDtEnd string //面试日期结束时间
	}
	serviceName := "serviceName"
	method := "serviceMethod"

	req := HReturnFeeRefreshReq{
		IntvDtSta: "2024-11-10",
		IntvDtEnd: "2024-11-10",
	}

	resp := new(interface{}) //返回接收
	err = ConsulRequest(consulClient, serviceName, method, req, resp)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(resp)
}
```

