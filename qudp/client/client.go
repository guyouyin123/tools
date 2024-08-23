package client

import (
	"fmt"
	"net"
)

func udpClient() {
	// 创建连接
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8080,
	})
	if err != nil {
		fmt.Println("连接失败!", err)
		return
	}
	defer socket.Close()
	// 发送给服务端数据
	sendData := []byte("hello server!")
	_, err = socket.Write(sendData)
	if err != nil {
		fmt.Println("发送数据失败!", err)
		return
	}

	// 接收客户端的数据
	data := make([]byte, 4096)
	read, remoteAddr, err := socket.ReadFromUDP(data)
	if err != nil {
		fmt.Println("读取数据失败!", err)
		return
	}
	fmt.Println(read, remoteAddr)
	fmt.Println("发送数据：", string(sendData))
}
