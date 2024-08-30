package server

import (
	"fmt"
	"log"
	"net"
)

// 处理用户请求
func HandleConn(conn net.Conn) {
	//函数调用完毕，自动关闭conn
	defer conn.Close()

	//获取客户端的网络地址信息
	addr := conn.RemoteAddr().String()
	log.Println(addr, " conncet sucessful")

	buf := make([]byte, 2048)

	for {
		//读取用户数据
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		log.Println("数据长度：", n)
		log.Println("接收数据：", string(buf[:n]))

		//if "exit" == string(buf[:n-1]) { //nc测试
		if "exit" == string(buf[:n-2]) { //自己写的客户端测试, 发送时，多了2个字符, "\r\n"
			log.Println(addr, " exit")
			return
		}

		sendAck := []byte("ack")
		n, err = conn.Write(sendAck)
		if err != nil {
			log.Println("err = ", err)
			return
		}
	}

}

func TcpServerRun() {
	//监听
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	defer listener.Close()

	//接收多个用户
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("err = ", err)
			return
		}
		//处理消息
		go HandleConn(conn)
	}

}
