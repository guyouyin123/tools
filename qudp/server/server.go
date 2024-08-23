package server

import (
	"fmt"
	"log"
	"net"
)

func UdpServer() {
	// 创建监听
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8080,
	})
	if err != nil {
		log.Println("listen err:", err)
		return
	}
	defer socket.Close()
	for {
		// 读取客户端传来的数据
		data := make([]byte, 4096)
		read, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			log.Println("读取数据失败!", err)
			continue
		}
		fmt.Println("收到数据：", string(data[:read]))
		//发送数据,告诉客户端已收到
		sendAck := []byte("ack")
		_, err = socket.WriteToUDP(sendAck, remoteAddr)
		if err != nil {
			log.Println("send data err:", err)
			return
		}
	}
}
