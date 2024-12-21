package qconsul

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"log"
	"strings"
)

// InitConn 初始化连接
func InitConn(address string) (*api.Client, error) {
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
		return nil, err
	}
	return client, nil
}

// RegisterServiceTest consul注册服务
func RegisterServiceTest(consulConn *api.Client) error {
	// 注册一个服务
	service := &api.AgentServiceRegistration{
		ID:      "helloWord-service",
		Name:    "helloWord",
		Port:    8500,
		Address: "127.0.0.1",
	}
	err := consulConn.Agent().ServiceRegister(service)
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
		return err
	}
	log.Println("Service registered successfully.")
	return nil
}

// ConsulRequest 发起调度
func ConsulRequest(consulClient *api.Client, serviceName, method string, req, resp interface{}) error {
	checks, _, err := consulClient.Health().Checks(serviceName, nil)
	if err != nil {
		log.Fatalf("Failed to find service %s: %v", serviceName, err)
		return err
	}
	passingAddress := []string{}
	for _, v := range checks {
		if v.Status == "passing" {
			s := strings.Split(v.ServiceID, "-")
			passingAddress = append(passingAddress, fmt.Sprintf("%s:%s", s[1], s[2]))
		} else {
			//注销
			err = consulClient.Agent().ServiceDeregister(v.ServiceID)
			if err != nil {
				log.Fatalf("Error deregistering service: %v", err)
				return err
			}
			log.Println("Service deregistered successfully.")
		}
	}
	if len(passingAddress) == 0 {
		log.Fatalf("Failed to find service %s: %v", serviceName, err)
		return err
	}

	serviceAddress := passingAddress[0]
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	ctx := context.Background()
	methodReq := fmt.Sprintf("%s.%s/%s", serviceName, serviceName, method)
	err = conn.Invoke(ctx, methodReq, req, resp)
	if err != nil {
		log.Fatalf("could not invoke: %v", err)
		return err
	}
	return nil
}
