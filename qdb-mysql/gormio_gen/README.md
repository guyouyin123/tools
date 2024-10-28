# db模型
```go
需要执行 main.go文件生成对应的dbmodel
```

```go
配置：./gen/conf/conf.go 配置文件

eg:

package conf

type DBConf struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	UserName string `json:"user"`
	PassWord string `json:"pass"`
	DBName   string `json:"db"`
}

var (
	//初始化的表名，主键名称
	IsTable = true //是否按需初始化表，如果是false就全量初始化所有表
	//按需加载输出话的表名
	TableNames = map[string]bool{
		"user":   true,
		"orders": true,
	}

	//db配置
	DbConf = DBConf{
		IP:       "127.0.0.1",
		Port:     3306,
		UserName: "root",
		PassWord: "123456",
		DBName:   "myself",
	}

	IsTemplate = true //是否开启模版函数
	/*
		模版配置--只增加了主键的查询方法--1.getId 2.getIds 3.GetIdsMap
		注意：不支持指针类型的健--字段类型允许为NULL生成的就为指针，可以FieldNullable=false，那么生成的就不是指针
	*/
	TemplateConf = map[string]map[string]string{
		"user": {
			"id":      "int32", //模版健和模版类型
			"id_card": "string",
		},
		"orders": {
			"id":       "int32",
			"order_id": "int32",
		},
	}
)
```



