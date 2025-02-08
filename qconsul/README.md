# qconsul-Consul--本文实现客户端服务发现，健康检查。rpc调用

是 HashiCorp 开发的一款工具，主要用于服务发现、配置管理和分布式系统的健康检查。它允许微服务之间进行通信，并提供一种简便的方法来管理服务的配置和状态。

**服务端官方下载链接**：https://developer.hashicorp.com/consul/install

后台启动命令：nohup ./consul agent -dev &



```go
package qconsul

import (
	"fmt"
	"testing"
	qconsul "github.com/guyouyin123/tools/qconsul"
)

func TestRunTest(t *testing.T) {
	Address := "127.0.0.1:8500"
	c := ConsulClient{IsDeregister: true, IndexType: 1}
	//初始化
	err := c.InitConn(Address)
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
	err = c.ConsulRequest(serviceName, method, req, resp)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(resp)
}

```

