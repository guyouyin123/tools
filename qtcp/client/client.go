package client

import (
	"log"
	"net"
)

func TcpClient() {
	//主动连接服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Println("err = ", err)
		return
	}

	defer conn.Close()

	//发送数据
	sendData := []byte("hello word!")
	conn.Write(sendData)
	log.Println("send Data:", string(sendData))

	//接收服务器回复的数据
	//切片缓冲
	buf := make([]byte, 1024)
	n, err := conn.Read(buf) //接收服务器的请求
	if err != nil {
		log.Println("conn.Read err = ", err)
		return
	}
	log.Println("收到回复内容：", string(buf[:n])) //打印接收到的内容, 转换为字符串再打印
}
