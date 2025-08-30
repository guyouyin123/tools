package qexcel

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"os"
	"testing"
	"time"
)

func TestWriteToXlsx(t *testing.T) {
	type User struct {
		Name  string `excel:"title=姓名;width=20;column=F"`
		Age   int    `excel:"title=年龄;width=50;column=A"`
		Email string `excel:"title=身份证;width=30;column=C"`
	}
	jeff := User{
		Name:  "jeff",
		Age:   18,
		Email: "12312312",
	}
	chary := User{
		Name:  "小明",
		Age:   20,
		Email: "xxxooo",
	}
	list := []interface{}{}
	list = append(list, jeff, chary)

	file, err := XlsxWriteV1(list, "Sheet1", "./test.xlsx", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(file)
}

func TestWriteToXlsxV2(t *testing.T) {
	type Info struct {
		Size    int        `json:"Size" excel:"title=面积;width=20;column=D"`
		Address string     `excel:"title=地址;width=20;column=E"`
		BeginDt *time.Time `json:"BeginDt" excel:"title=开始日期;width=20;column=D"` // 开始日期

	}
	type Like struct {
		Desc string `excel:"title=爱好;width=40;column=F"`
	}
	type Friend struct {
		Friend string `excel:"title=朋友;width=40;column=G"`
	}
	type User struct {
		ID     int    `json:"-"`
		Name   string `json:"-" excel:"title=姓名;width=20;column=A"`
		Age    int    `excel:"title=年龄;width=40;column=B"`
		Email  string `excel:"title=身份证;width=30;column=C"`
		House  []*Info
		Like   []*Like
		Friend []*Friend
	}

	now := time.Now()
	h1 := &Info{
		Size:    100,
		Address: "上海市浦东新区东方明珠",
		BeginDt: &now,
	}
	h2 := &Info{
		Size:    90,
		Address: "上海市xxxooo",
	}
	h3 := &Info{
		Size:    80,
		Address: "北京市xxxooo",
	}
	f1 := &Like{
		Desc: "打球",
	}
	f2 := &Like{
		Desc: "泡澡",
	}
	jeffList := make([]*Info, 0)
	jeffList = append(jeffList, h1)
	jeffList = append(jeffList, h2)
	jeffList = append(jeffList, h3)

	frList := make([]*Like, 0)
	frList = append(frList, f1)
	frList = append(frList, f2)

	fList := make([]*Friend, 0)
	for i := 0; i < 10; i++ {
		p1 := Friend{
			Friend: fmt.Sprintf("p%d", i),
		}
		fList = append(fList, &p1)
	}

	jeff := User{
		Name:  "jeff",
		Age:   18,
		Email: "12312312",
		House: jeffList,
		Like:  frList,
	}
	chary := User{
		Name:   "小明",
		Age:    20,
		Email:  "xxxooo",
		Like:   frList,
		Friend: fList,
	}
	list := []interface{}{}
	list = append(list, jeff)
	list = append(list, chary)

	file, err := XlsxWriteV2(list, "Sheet1", "./test1.xlsx", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(file)
}

func TestExportCsv(t *testing.T) {
	var csvData [][]string
	csvTitle := []string{"姓名", "身份证号码", "手机号码", "地址"}
	data1 := []string{"小张", "123456", "123321", `xxoo`}
	data2 := []string{"小明", "654321", "123333", `ooxx`}
	csvData = append(csvData, data1)
	csvData = append(csvData, data2)
	data, err := CsvWrite(csvTitle, csvData)
	if err != nil {
		t.Fatal(err)
	}
	temp, err := os.Create("./test.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer temp.Close()
	temp.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	temp.Write(data)
}

func TestWriteToXlsxV2_1(t *testing.T) {
	type User struct {
		UserId   int    `excel:"title=用户id;width=20;column=B"`
		UserName string `excel:"title=用户名称;width=20;column=C"`
	}
	type User2 struct {
		UserId2   int    `excel:"title=用户id;width=20;column=D"`
		UserName2 string `excel:"title=用户名称;width=20;column=E"`
	}

	type UserBase struct {
		BossId int64 `excel:"title=大佬;width=20;column=A"`
		Data   []*User
		Data2  []*User2
	}
	a := &User{
		UserId:   1,
		UserName: "小明",
	}
	aList := make([]*User, 0)
	aList = append(aList, a)
	aList = append(aList, a)

	b := &UserBase{
		BossId: 100,
		Data:   aList,
	}
	exportList := make([]*UserBase, 0, 10)
	exportList = append(exportList, b)

	list := make([]interface{}, 0)
	for _, v := range exportList {
		list = append(list, v)
	}
	sheetName := "Sheet1"
	_, err := XlsxWriteV2(list, sheetName, "./test2.xlsx", true)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func TestWriteToXlsxV3(t *testing.T) {
	type Desc struct {
		Url string `excel:"title=url;width=20;column=E"`
	}

	type like struct {
		Like string `excel:"title=爱好;width=20;column=D"`
		Desc []*Desc
	}

	type Class struct {
		Class string `excel:"title=分类;width=20;column=C"`
		Like  []*like
	}
	type User struct {
		Name  string `excel:"title=姓名;width=20;column=F"`
		Age   int    `excel:"title=年龄;width=20;column=B"`
		Type  int    `excel:"title=类型;width=20;column=A;enum={\"1\":\"老师\",\"2\":\"学生\",\"0\":\"未知\"}"`
		Class []*Class
	}

	l1 := like{
		Like: "打球",
		Desc: []*Desc{
			{Url: "www.baidu1.com"},
			{Url: "www.baidu2.com"},
		},
	}
	l2 := like{
		Like: "跑步",
		Desc: []*Desc{
			{Url: "www.baidu3.com"},
			{Url: "www.baidu4.com"},
		},
	}
	l3 := like{
		Like: "三国",
		Desc: []*Desc{
			{Url: "www.baidu5.com"},
			{Url: "www.baidu6.com"},
		},
	}
	l4 := like{
		Like: "水浒",
		Desc: []*Desc{
			{Url: "www.baidu7.com"},
			{Url: "www.baidu8.com"},
		},
	}

	like1 := Class{
		Class: "球类",
		Like:  []*like{&l1, &l2},
	}
	like2 := Class{
		Class: "娱乐类",
		Like:  []*like{&l3, &l4},
	}

	user1 := User{
		Name:  "Jeff",
		Age:   18,
		Class: []*Class{&like1, &like2},
		Type:  99,
	}

	sheetName := "Sheet1"
	_, err := XlsxWriteV3(nil, user1, sheetName, "./user2.xlsx", true)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Test1WriteToXlsxV3(t *testing.T) {
	/*
			wrap_text:true //自动换行
			vertical:"top" //垂直对齐方式
			horizontal:"center" //居中对齐方式
			indent:1 //缩进
			shrink_to_fit:false //不缩小字体填充
			text_rotation:0 //文本旋转角度
		"font": {
		    "color": "#FF0000" //字体颜色
		  },
		  "fill": {
		    "type": "pattern",
		    "color": ["#FFFF00"], //背景颜色
		    "pattern": 1
		  }
	*/
	type like struct {
		Like string `excel:"title=爱好;width=20;column=C;style={\"alignment\":{\"horizontal\":\"center\",\"vertical\":\"center\",\"wrap_text\":true}}"`
	}
	type User struct {
		Name          string `excel:"title=姓名;width=20;column=A;style={\"alignment\":{\"horizontal\":\"center\",\"text_rotation\":45},\"font\":{\"color\":\"#FF0000\"},\"fill\":{\"type\":\"pattern\",\"color\":[\"#FFFF00\"],\"pattern\":1}}"`
		Age           int
		SettlementTyp int8 `excel:"title=模式;width=10;column=F;enum={\"0\":\"未知\",\"1\":\"ZX模式\",\"2\":\"Z模式\",\"3\":\"ZA模式\",\"4\":\"Z-B模式\",\"5\":\"ZX-B模式\",\"6\":\"ZX-A模式\",\"7\":\"Z-D模式\",\"8\":\"ZX-D模式\"}"`
		Like          []*like
	}
	list := []*User{
		{
			Name:          "张三",
			Age:           18,
			SettlementTyp: 1,
			Like: []*like{
				{
					Like: "吃饭很长的描述测试自动换行，吃饭很长的描述测试自动换行，吃饭很长的描述测试自动换行",
				},
				{
					Like: "睡觉",
				},
			},
		},
		{
			Name:          "李四",
			Age:           19,
			SettlementTyp: 2,
			Like: []*like{
				{
					Like: "123",
				},
				{
					Like: "456",
				},
			},
		},
	}

	sheetName := "测试的excel"
	_, err := XlsxWriteV3(nil, &list, sheetName, "./userv3.xlsx", true)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestDemo1ToXlsxV3(t *testing.T) {
	type AbnormalInfo struct {
		AbnormalClass         int64  `excel:"title=异常类型;width=20;column=F;"`
		AbnormalType          int64  //异常指标
		AbnormalName          string `excel:"title=异常指标;width=20;column=G;"`
		AbnormalValue         string `excel:"title=值;width=20;column=H;"`
		Money                 int64  `excel:"title=金额;width=20;column=I;"`
		Status                int64  `excel:"title=预开状态;width=20;column=J;"`
		WdServiceStatisticsID int64  //罚单id
	}
	type FineForBrokerRecordList struct {
		WdServiceStatisticsId int64  //罚单id
		DayDt                 string `excel:"title=日期;width=20;column=A;"`
		DayDtYear             string //罚单所属日期
		MiddleAreaName        string `excel:"title=中小区;width=20;column=B;"`
		AreaName              string `excel:"title=小区;width=20;column=C;"`
		StoreName             string `excel:"title=门店;width=20;column=D;"`
		BrokerUserId          int64  //经纪人id
		BrokerName            string `excel:"title=经纪人姓名;width=20;column=E;"`
		BrokerReaName         string //经纪人姓名
		Msg                   string //理由
		Replay                string `excel:"title=数据复盘;width=20;column=K;"`
		Reason                string `excel:"title=原因;width=20;column=L;"`
		Scheme                string `excel:"title=方案;width=20;column=M;"`
		Abnormal              []*AbnormalInfo
	}

	a := AbnormalInfo{
		AbnormalClass:         1,
		AbnormalType:          1,
		AbnormalName:          "1",
		AbnormalValue:         "1",
		Money:                 1,
		Status:                1,
		WdServiceStatisticsID: 1,
	}

	b := FineForBrokerRecordList{
		WdServiceStatisticsId: 1,
		DayDt:                 "1",
		DayDtYear:             "1",
		MiddleAreaName:        "1",
		AreaName:              "1",
		StoreName:             "1",
		BrokerUserId:          1,
		BrokerName:            "1",
		BrokerReaName:         "1",
		Abnormal:              []*AbnormalInfo{&a, &a},
		Msg:                   "1",
	}
	list := []*FineForBrokerRecordList{&b, &b, &b, &b, &b, &b, &b, &b}

	sheetName := "测试的excel"
	_, err := XlsxWriteV3(nil, &list, sheetName, "./record.xlsx", true)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestDemo1(t *testing.T) {
	type Friend struct {
		FrName string `excel:"title=女朋友们;width=20;column=C;"`
	}
	type User struct {
		Name       string `excel:"title=姓名;width=20.5;column=A;IsMerge=true"`
		Age        int64  `excel:"title=年龄;width=20;column=B;IsMerge=true"`
		IdCard     int64  `excel:"title=身份证;width=20;column=D;IsMerge=true"`
		FriendList []*Friend
	}

	f1 := Friend{
		FrName: "杨贵妃",
	}
	f2 := Friend{
		FrName: "三上",
	}
	f3 := Friend{
		FrName: "小优",
	}
	Jeff := User{
		Name:       "Jeff",
		Age:        18,
		FriendList: []*Friend{&f1, &f2, &f3},
	}
	Jeff2 := User{
		Name:       "Jeff2",
		Age:        20,
		FriendList: []*Friend{&f1, &f2},
	}
	list := []*User{&Jeff, &Jeff2}

	sheetName := "测试的excel"
	_, err := XlsxWriteV3(nil, &list, sheetName, "./record.xlsx", true)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestDemo2(t *testing.T) {
	type AbnormalInfo struct {
		//WdServiceStatisticsID int64 `excel:"title=ID;width=12.28;column=O;"` //罚单id
		WdServiceStatisticsID int64  //罚单id
		AbnormalClass         int64  `excel:"title=异常类型;width=12.28;column=F;"` //异常类型 1工作量 2服务质量
		AbnormalType          int64  //异常指标
		AbnormalName          string `excel:"title=异常指标;width=21.8;column=G;"` //异常名称
		FineValue             string `excel:"title=标准;width=7.74;column=H;"`
		AbnormalValue         string `excel:"title=值;width=7.74;column=I;"`    //异常值
		Money                 int64  `excel:"title=金额;width=7.74;column=J;"`   //金额(单位分)
		Status                int64  `excel:"title=预开状态;width=7.74;column=K;"` //罚单状态 -1不满足 1未开 2已撤销 3已开 999 全部
	}
	type FineForBrokerRecordList struct {
		WdServiceStatisticsId int64  //罚单id
		DayDt                 string `excel:"title=日期;width=7.28;column=A;IsMerge=true;style={\"alignment\":{\"horizontal\":\"left\",\"vertical\":\"top\",\"wrap_text\":true}}"`
		DayDtYear             string //罚单所属日期
		MiddleAreaName        string `excel:"title=中小区;width=7.28;column=B;IsMerge=true;style={\"alignment\":{\"horizontal\":\"left\",\"vertical\":\"top\",\"wrap_text\":true}}"`
		AreaName              string `excel:"title=小区;width=7.28;column=C;IsMerge=true;style={\"alignment\":{\"horizontal\":\"left\",\"vertical\":\"top\",\"wrap_text\":true}}"`
		StoreName             string `excel:"title=门店;width=12.28;column=D;IsMerge=true;style={\"alignment\":{\"horizontal\":\"left\",\"vertical\":\"top\",\"wrap_text\":true}}"`
		BrokerUserId          int64  //经纪人id
		BrokerName            string `excel:"title=经纪人;width=12.28;column=E;IsMerge=true;style={\"alignment\":{\"horizontal\":\"left\",\"vertical\":\"top\",\"wrap_text\":true}}"`
		BrokerReaName         string //经纪人姓名
		Msg                   string //理由
		Replay                string `excel:"title=数据复盘;width=23.8;column=L;IsMerge=true;style={\"alignment\":{\"horizontal\":\"left\",\"vertical\":\"top\",\"wrap_text\":true}}"`
		Reason                string `excel:"title=原因;width=23.8;column=M;IsMerge=true;style={\"alignment\":{\"horizontal\":\"left\",\"vertical\":\"top\",\"wrap_text\":true}}"`
		Scheme                string `excel:"title=方案;width=23.8;column=N;IsMerge=true;style={\"alignment\":{\"horizontal\":\"left\",\"vertical\":\"top\",\"wrap_text\":true}}"`
		AppealStatus          int64  //申诉状态 0未申诉 1已申诉
		Abnormal              []*AbnormalInfo
	}

	list := []*FineForBrokerRecordList{}
	jsoniter.Unmarshal([]byte(dataStr), &list)

	sheetName := "测试excel"
	_, err := XlsxWriteV3(nil, &list, sheetName, "./test.xlsx", true)
	if err != nil {
		fmt.Println(err)
		return
	}
}

const dataStr = `[{}]`
