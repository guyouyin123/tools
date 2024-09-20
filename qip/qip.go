package qip

import (
	"fmt"
	"github.com/guyouyin123/tools/qhttp"
	jsoniter "github.com/json-iterator/go"
)

// LocationInfo 存储 IP 归属地信息
type LocationInfo struct {
	Query       string  `json:"query"`       //查询的IP
	Status      string  `json:"status"`      //状态success
	Country     string  `json:"country"`     //国家
	CountryCode string  `json:"countryCode"` //国家编码
	Region      string  `json:"region"`      //区域
	RegionName  string  `json:"regionName"`  //区域名称
	City        string  `json:"city"`        //城市
	Zip         string  `json:"zip"`         //96521
	Lat         float64 `json:"lat"`         //纬度
	Lon         float64 `json:"lon"`         //经度
	Timezone    string  `json:"timezone"`    //时区
	ISP         string  `json:"isp"`         //提供者
	Org         string  `json:"org"`         //提供组织
	AS          string  `json:"AS"`          //提供公司
}

// GetIPLocation 查询 IP 归属地
func GetIPLocation(ip string) (*LocationInfo, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := qhttp.Get(url, nil, nil)
	if err != nil {
		return nil, err
	}
	location := &LocationInfo{}
	_ = jsoniter.Unmarshal(resp, &location)
	return location, nil
}
