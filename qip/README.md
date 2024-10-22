# ip地址解析

```go
使用：url := fmt.Sprintf("http://ip-api.com/json/%s", ip) //本文使用方式--测试更精准
使用：url := fmt.Sprintf("https://ipinfo.io/%s/json", ip) //不是很精准
```



```go
func TestGetIPLocation(t *testing.T) {
   ip := "47.100.254.82" // 替换为要查询的 IP 地址
   location, err := GetIPLocation(ip)
   if err != nil {
      fmt.Println("Error:", err)
      return
   }
   fmt.Println(location)
}
```