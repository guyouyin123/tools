package main

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func main() {
	config := func(conf *neo4j.Config) {
		//conf.Encrypted = false                         // 启用 TLS
		//conf.TLSCertFile = "/path/to/certificate.crt" // 设置 TLS 证书文件路径
		//conf.TLSKeyFile = "/path/to/private.key"      // 设置 TLS 私钥文件路径
		//conf.TLSCAFile = "/path/to/ca.crt"            // 设置 TLS CA 证书文件路径
	}

	// 创建驱动程序
	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "123456", ""), config)
	if err != nil {
		// 处理错误
	}
	defer driver.Close()
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	cyhperStr := `CREATE (n:city2 {name: $name})` //创建节点
	//建立关系
	cyhperStr2 := `MATCH (n:city2 {name: $name}),(k:city {name: "苏州"}) CREATE (n)-[:属于]->(k)`
	nameList := []string{"昆山"}
	//查询节点存在否
	cyhperStr3 := `MATCH (node:city2)
WHERE node.name = $name
RETURN node`
	for _, name := range nameList {
		//查询节点
		_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run(cyhperStr3, map[string]interface{}{
				"name": name,
			})
			if err != nil {
				return false, err
			}

			//存在则不创建，不存在就创建
			if result.Next() {
				return result, nil
			}

			//创建节点
			result, err = tx.Run(cyhperStr, map[string]interface{}{
				"name": name,
			})
			if err != nil {
				return false, err
			}

			//创建关系
			result, err = tx.Run(cyhperStr2, map[string]interface{}{
				"name": name,
			})
			if err != nil {
				return nil, err
			}
			return result, nil
		})

		if err != nil {
			// 处理错误
			fmt.Println("err:", err)
			return
		}
	}

}
