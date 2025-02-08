package qconsul

import (
	"fmt"
	"testing"
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
