package qconsul

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"log"
	"strings"
	"time"
)

type ConsulClient struct {
	client     *api.Client                   //consul 客户端
	connRpcMap map[string][]*grpc.ClientConn //Rpc连接池
}

// InitConn 初始化连接
func (conn *ConsulClient) InitConn(address string) error {
	// 创建一个新的 Consul 客户端
	config := &api.Config{
		Address:    address, // Consul 服务器的地址和端口号
		Scheme:     "",      // 协议可选，默认为 http 可设置为https
		PathPrefix: "",      //自动添加 api的 前缀
		Datacenter: "",      //Consul 数据中心
		Transport:  nil,
		HttpClient: nil,
		HttpAuth:   nil,
		WaitTime:   0,  // 超时时间
		Token:      "", //Consul API 的 ACL 令牌
		TokenFile:  "", //访问 Consul API 的 ACL 令牌
		Namespace:  "", //命名空间隔离，与k8s名称空间相似
		Partition:  "", //分区隔离，可同一个Consul不同分区，例如测试和产线
		TLSConfig:  api.TLSConfig{},
	}
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
		return err
	}
	conn.client = client
	conn.connRpcMap = make(map[string][]*grpc.ClientConn, 0)
	go conn.CheckRpcConn()
	return nil
}

// CheckRpcConn 健康检查Rpc连接池
func (conn *ConsulClient) CheckRpcConn() {
	timer := time.NewTicker(60 * time.Second)
	for {
		<-timer.C
		serviceNameList := make([]string, 0)
		for serviceName, _ := range conn.connRpcMap {
			serviceNameList = append(serviceNameList, serviceName)
		}
		for _, serviceName := range serviceNameList {
			_, err := conn.GetConnList(serviceName)
			if err != nil {
				conn.connRpcMap[serviceName] = nil
				log.Fatalf("Failed to find service %s: %v", serviceName, err)
			}
		}
	}
}

// RegisterServiceTest consul注册服务
func (conn *ConsulClient) RegisterServiceTest() error {
	// 注册一个服务
	service := &api.AgentServiceRegistration{
		ID:      "helloWord-service",
		Name:    "helloWord",
		Port:    8500,
		Address: "127.0.0.1",
	}
	err := conn.client.Agent().ServiceRegister(service)
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
		return err
	}
	log.Println("Service registered successfully.")
	return nil
}

// ConsulRequest 发起调度
func (conn *ConsulClient) ConsulRequest(serviceName, method string, req, resp interface{}) error {
	var err error
	connRpcList, ok := conn.connRpcMap[serviceName]
	if !ok {
		connRpcList, err = conn.GetConnList(serviceName)
		if err != nil {
			return err
		}
	}
	connRpc := connRpcList[0] //TODO 调度算法,轮训，随机，权重等
	methodReq := fmt.Sprintf("%s.%s/%s", serviceName, serviceName, method)
	err = connRpc.Invoke(context.Background(), methodReq, req, resp)
	if err != nil {
		log.Fatalf("could not invoke: %v", err)
		return err
	}
	return nil
}

func (conn *ConsulClient) GetConnList(serviceName string) ([]*grpc.ClientConn, error) {
	connList := make([]*grpc.ClientConn, 0)
	checks, _, err := conn.client.Health().Checks(serviceName, nil)
	if err != nil {
		log.Fatalf("Failed to find service %s: %v", serviceName, err)
		return nil, err
	}
	passingAddress := []string{}
	for _, v := range checks {
		if v.Status == "passing" {
			s := strings.Split(v.ServiceID, "-")
			passingAddress = append(passingAddress, fmt.Sprintf("%s:%s", s[1], s[2]))
		} else {
			//注销
			err = conn.client.Agent().ServiceDeregister(v.ServiceID)
			if err != nil {
				log.Fatalf("Error deregistering service: %v", err)
				return nil, err
			}
			log.Println("Service deregistered successfully.")
		}
	}
	if len(passingAddress) == 0 {
		log.Fatalf("Failed to find service %s: %v", serviceName, err)
		return nil, err
	}
	for _, serviceAddress := range passingAddress {
		connRpc, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
			return nil, err
		}
		connList = append(connList, connRpc)
	}
	conn.connRpcMap[serviceName] = connList
	return connList, nil
}
