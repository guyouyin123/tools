package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"net/http"
	"time"
	"unsafe"
)

type WsServer struct {
	clientMap map[string]*websocket.Conn
	MaxCount  int //最大连接数
	Timer     int //心跳间隔，单位秒
}

// 发送消息结构体
type RespMsg struct {
	Status int         `json:"status"` //状态码
	Type   int         `json:"type"`   //消息类型
	Data   interface{} `json:"data"`   //消息体
	Time   string      `json:"time"`   //消息时间
}

var wsServer *WsServer

// 初始化
func InitWs() *WsServer {
	return &WsServer{
		clientMap: make(map[string]*websocket.Conn, 10),
		MaxCount:  3,
		Timer:     5,
	}
}

// 读取数据
func (server *WsServer) readLoop(ws *websocket.Conn) (chanU chan []byte) {
	messageChan := make(chan []byte) //每个ws连接独立的chan
	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				server.close(ws)
				return
			}
			messageChan <- message
		}
	}()
	return messageChan
}

// 向客户端发送数据
func (server *WsServer) sendLoop(ws *websocket.Conn, ack interface{}, messageType int) error {
	s, _ := json.Marshal(ack)
	msg := *(*[]byte)(unsafe.Pointer(&s))
	//服务端发送消息到客户端websocket
	err := ws.WriteMessage(messageType, msg) //messageType 1.文本消息 2.二进制消息 8.关闭消息 9.ping消息 10.pong消息
	if err != nil {
		return err
	}
	fmt.Println("发送消息：", string(s))
	return nil
}

// 广播群发消息--帧同步
func (server *WsServer) broadcast(ack interface{}) {
	//写入ws数据
	for _, ws := range server.clientMap {
		resp := RespMsg{
			Status: 200,
			Type:   1,
			Time:   time.Now().Format("2006-01-02 15:04:05"),
			Data:   &ack,
		}
		if err := server.sendLoop(ws, resp, 1); err != nil {
			log.Error(err.Error())
			return
		}
	}
}

// 新建连接
func (server *WsServer) NewClientWs(c *gin.Context) (*websocket.Conn, error) {
	//连接数量限制
	if len(server.clientMap) >= server.MaxCount {
		fmt.Println("连接数超限：", wsServer.MaxCount)
		return nil, errors.New("连接数超限")
	}

	// 设置websocket
	// CheckOrigin防止跨站点的请求伪造
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}
	go server.ping(ws) //心跳
	key := fmt.Sprintf("%p", ws)
	wsServer.clientMap[key] = ws //加入连接map
	return ws, nil
}

// 关闭
func (server *WsServer) close(ws *websocket.Conn) {
	ws.Close()
	//移除链接池
	key := fmt.Sprintf("%p", ws)
	delete(server.clientMap, key)
	fmt.Println("=======删除连接，剩余：", len(wsServer.clientMap))
}

// 心跳
func (server *WsServer) ping(ws *websocket.Conn) {
	for {
		time.Sleep(time.Second * time.Duration(server.Timer))
		ack := "心跳包"
		resp := RespMsg{
			Status: 200,
			Type:   2,
			Time:   time.Now().Format("2006-01-02 15:04:05"),
			Data:   &ack,
		}
		//设置为 9，心跳包，客户端输出将不可见
		if err := server.sendLoop(ws, resp, 9); err != nil {
			fmt.Println("心跳停止，停止心跳包发送")
			return
		}
	}
}

// websocket实现--业务层
func pk(c *gin.Context) {
	//新建ws连接
	ws, err := wsServer.NewClientWs(c)
	if err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
		return
	}

	//监听读取数据
	messageChan := wsServer.readLoop(ws)
	//消息处理
	for {
		select {
		case message := <-messageChan:
			//读取客户端发送来到消息
			if message != nil {
				fmt.Println("服务端收到消息:", string(message))

				type s struct {
					Msg string `json:"msg"` //广播内容
				}
				ack := s{Msg: "广播内容"}
				wsServer.broadcast(&ack) //广播群发消息
			}
		}
	}
}

func Run() {
	r := gin.Default()
	wsServer = InitWs()

	r.GET("/ws", pk)
	r.Run(":8080")

}
