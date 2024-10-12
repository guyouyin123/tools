package qconsul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
)

func ConsulClient() {
	// 创建一个新的 Consul 客户端
	config := api.DefaultConfig()
	config.Datacenter = "zxx"
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	// 注册一个服务
	service := &api.AgentServiceRegistration{
		ID:      "helloWord-service",
		Name:    "helloWord",
		Port:    8500,
		Address: "127.0.0.1",
	}

	err = client.Agent().ServiceRegister(service)
	if err != nil {
		log.Fatalf("Error registering service: %v", err)
	}
	log.Println("Service registered successfully.")

	// 查询服务
	checks, _, err := client.Health().Checks("helloWord", &api.QueryOptions{
		Datacenter: "zxx",
	})
	if err != nil {
		log.Fatalf("Error retrieving services: %v", err)
	}
	// 打印健康检查信息
	for _, check := range checks {
		log.Printf("Check ID: %s, Name: %s, Status: %s, Service ID: %s",
			check.CheckID, check.Name, check.Status, check.ServiceID)
	}

	fmt.Println(checks)
	// 注销服务
	//err = client.Agent().ServiceDeregister(service.ID)
	//if err != nil {
	//	log.Fatalf("Error deregistering service: %v", err)
	//}
	//log.Println("Service deregistered successfully.")
}
