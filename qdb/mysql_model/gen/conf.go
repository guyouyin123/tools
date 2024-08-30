package main

type DBConf struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	UserName string `json:"user"`
	PassWord string `json:"pass"`
	DBName   string `json:"db"`
}

// 本项目需要用到的表名配置这里，避免无用的表多余加载。--可以取消，加载全部表
var TableNames = map[string]bool{
	"user": true,
}

var dbConf = DBConf{
	IP:       "127.0.0.1",
	Port:     3306,
	UserName: "root",
	PassWord: "123456",
	DBName:   "test",
}
