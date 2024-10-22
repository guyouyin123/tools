# 压测工具

#下载地址文档 https://github.com/link1st/go-stress-testing

Usage of ./go-stress-testing-mac:
-c uint
并发数 (default 1)
-n uint
请求数(单个并发/协程) (default 1)
-u string
压测地址
-d string
调试模式 (default "false")
-http2
是否开http2.0
-k	是否开启长连接
-m int
单个host最大连接数 (default 1)
-H value
自定义头信息传递给服务器 示例:-H 'Content-Type: application/json'
-data string
HTTP POST方式传送数据
-v string
验证方法 http 支持:statusCode、json webSocket支持:json
-p string
curl文件路径