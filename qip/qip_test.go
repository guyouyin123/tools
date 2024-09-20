package qip

import (
	"fmt"
	"testing"
)

func TestGetIPLocation(t *testing.T) {
	//ip := "47.100.254.82" // 替换为要查询的 IP 地址
	//ip := "180.88.111.187" // 替换为要查询的 IP 地址
	ip := "1.208.112.21" // 替换为要查询的 IP 地址
	location, err := GetIPLocation(ip)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(location)
}
