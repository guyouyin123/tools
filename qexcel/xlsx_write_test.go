package qexcel

import (
	"fmt"
	"strings"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestWriteToXlsx(t *testing.T) {
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
	_, err := XlsxWrite(nil, user1, sheetName, "./user2.xlsx", true)
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
	_, err := XlsxWrite(nil, &list, sheetName, "./userv3.xlsx", true)
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
		DayDt                 string `excel:"title=日期;width=20;column=A;IsMerge=true"`
		DayDtYear             string //罚单所属日期
		Abnormal              []*AbnormalInfo
		MiddleAreaName        string `excel:"title=中小区;width=20;column=B;IsMerge=true"`
		AreaName              string `excel:"title=小区;width=20;column=C;IsMerge=true"`
		StoreName             string `excel:"title=门店;width=20;column=D;IsMerge=true"`
		BrokerUserId          int64  //经纪人id
		BrokerName            string `excel:"title=经纪人姓名;width=20;column=E;IsMerge=true"`
		BrokerReaName         string //经纪人姓名
		Msg                   string //理由
		Replay                string `excel:"title=数据复盘;width=20;column=K;IsMerge=true"`
		Reason                string `excel:"title=原因;width=20;column=L;IsMerge=true"`
		Scheme                string `excel:"title=方案;width=20;column=M;IsMerge=true"`
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
	_, err := XlsxWrite(nil, &list, sheetName, "./record.xlsx", true)
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
	_, err := XlsxWrite(nil, &list, sheetName, "./record.xlsx", true)
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
	_, err := XlsxWrite(nil, &list, sheetName, "./test.xlsx", true)
	if err != nil {
		fmt.Println(err)
		return
	}
}

const dataStr = `[{}]`

func TestComplexNestedSlice(t *testing.T) {
	// 模拟3层嵌套: User -> []Order -> []Item
	type Item struct {
		ItemName string `excel:"title=商品名;width=15;column=C"`
		Price    int    `excel:"title=价格;width=10;column=D"`
	}
	type Order struct {
		OrderId string  `excel:"title=订单号;width=20;column=B;IsMerge=true"`
		Items   []*Item // 嵌套切片
	}
	type User struct {
		UserName string   `excel:"title=用户名;column=A;IsMerge=true"`
		Orders   []*Order // 嵌套切片
	}

	// 构造数据
	// User1 has 2 orders
	// Order1 has 2 items
	// Order2 has 1 item
	// Total rows for User1 = 2 + 1 = 3 rows

	//item1 := &Item{ItemName: "Apple", Price: 10}
	//item2 := &Item{ItemName: "Banana", Price: 5}
	//item3 := &Item{ItemName: "Orange", Price: 8}

	//order1 := &Order{
	//	OrderId: "ORD001",
	//	Items:   []*Item{item1, item2},
	//}
	//order2 := &Order{
	//	OrderId: "ORD002",
	//	Items:   []*Item{item3},
	//}

	//user1 := &User{
	//	UserName: "Alice",
	//	Orders:   []*Order{order1, order2},
	//}

	// User2 has 1 order with 0 items (should occupy 1 row)
	order3 := &Order{
		OrderId: "ORD003",
		Items:   []*Item{}, // Empty items
	}
	//user2 := &User{
	//	UserName: "Bob",
	//	Orders:   []*Order{order3},
	//}

	//list := []*User{user1, user2}

	_, err := XlsxWrite(nil, &order3, "DeepNested", "/Users/jeff/Desktop/deep_nested.xlsx", true)
	if err != nil {
		t.Fatalf("Failed to write deep nested struct: %v", err)
	}
}

// TestPointerSliceNestedStructPanic 回归: []*Parent + 非切片嵌套结构体
//
// 曾经的 bug: 传入 []*Parent(元素是指针)时 dataList[0] 为 *Parent, _tagHandle 的 baseVa 是 Ptr;
// 若 Parent 含一个"非切片"的嵌套结构体字段(下面的 Addr), 会走 `case reflect.Struct` 分支,
// 该分支直接 baseVa.Field(i) 未解引用指针 -> panic "call of reflect.Value.Field on ptr Value"。
// (对比: 嵌套"切片"走 case reflect.Slice 分支, 那里有 baseVa.Elem() 解引用, 故一直没事。)
//
// 修复后此形状应正常导出, 这里断言 err == nil。
func TestPointerSliceNestedStructPanic(t *testing.T) {
	type Addr struct {
		City string `excel:"title=城市;width=15;column=B"`
	}
	type Person struct {
		Name string `excel:"title=姓名;width=15;column=A"`
		Addr Addr   // 关键: 非切片的嵌套结构体
	}

	list := []*Person{
		{Name: "Tom", Addr: Addr{City: "北京"}},
	}

	_, err := XlsxWrite(nil, &list, "Sheet1", "./ptr_nested.xlsx", true)
	if err != nil {
		t.Fatalf("修复后 []*Parent+非切片嵌套结构体 应正常导出, 却报错: %v", err)
	}
}

// TestSameNameFieldNoCollision 回归(B): 不同嵌套结构体里的同名字段(Status)不应互相覆盖列
//
// 修复前 tagMap 以 field.Name 为 key, Inner.Status 会覆盖 Outer.Status,
// 导致两个 Status 都写到 B 列、A 列为空。修复后各归各列: A2=OUT, B2=IN。
func TestSameNameFieldNoCollision(t *testing.T) {
	type Inner struct {
		Status string `excel:"title=内层状态;width=15;column=B"`
	}
	type Outer struct {
		Status string `excel:"title=外层状态;width=15;column=A"` // 与 Inner.Status 同名不同列
		Inner  Inner
	}

	data := Outer{Status: "OUT", Inner: Inner{Status: "IN"}}
	f, err := XlsxWrite(nil, data, "Sheet1", "", false)
	if err != nil {
		t.Fatalf("导出失败: %v", err)
	}
	if got, _ := f.GetCellValue("Sheet1", "A2"); got != "OUT" {
		t.Fatalf("A2 期望 OUT, 实得 %q(外层 Status 被覆盖了)", got)
	}
	if got, _ := f.GetCellValue("Sheet1", "B2"); got != "IN" {
		t.Fatalf("B2 期望 IN, 实得 %q", got)
	}
}

// TestSinglePointerStruct 回归(F): 传入单个 *struct(指向结构体的指针)不应输出空文件
//
// 修复前 XlsxWrite 的 Ptr 分支只处理 *[]T, 指向单结构体时 dataList 为空 -> 直接返回空表。
// 修复后应正常写入: A2=Jeff, B2=18。
func TestSinglePointerStruct(t *testing.T) {
	type User struct {
		Name string `excel:"title=姓名;width=15;column=A"`
		Age  int    `excel:"title=年龄;width=10;column=B"`
	}

	f, err := XlsxWrite(nil, &User{Name: "Jeff", Age: 18}, "Sheet1", "", false)
	if err != nil {
		t.Fatalf("导出失败: %v", err)
	}
	if got, _ := f.GetCellValue("Sheet1", "A2"); got != "Jeff" {
		t.Fatalf("A2 期望 Jeff, 实得 %q(单个 *struct 输出了空表)", got)
	}
	if got, _ := f.GetCellValue("Sheet1", "B2"); got != "18" {
		t.Fatalf("B2 期望 18, 实得 %q", got)
	}

	// 二级指针 **struct(如 order:=&Item{}; XlsxWrite(&order,...)): 只有标题、无值内容 的回归。
	// 根因是 indirect 只解一层指针, **struct 解一层后是 *struct 非 struct -> renderItem 跳过。
	pp := &User{Name: "Bob", Age: 20}
	f2, err := XlsxWrite(nil, &pp, "Sheet1", "", false)
	if err != nil {
		t.Fatalf("二级指针导出失败: %v", err)
	}
	if got, _ := f2.GetCellValue("Sheet1", "A2"); got != "Bob" {
		t.Fatalf("二级指针 A2 期望 Bob, 实得 %q(只出标题没内容)", got)
	}
	if got, _ := f2.GetCellValue("Sheet1", "B2"); got != "20" {
		t.Fatalf("二级指针 B2 期望 20, 实得 %q", got)
	}
}

// TestMissingColumnError 回归(C): 字段写了 excel tag 但漏 column= 应直接报错而非静默写不进去
func TestMissingColumnError(t *testing.T) {
	type Bad struct {
		Name string `excel:"title=姓名;width=10"` // 故意缺 column
	}
	_, err := XlsxWrite(nil, Bad{Name: "x"}, "Sheet1", "", false)
	if err == nil {
		t.Fatalf("缺少 column 时应报错, 实际却成功")
	}
	if !strings.Contains(err.Error(), "column") {
		t.Fatalf("错误信息应指出缺少 column, 实得: %v", err)
	}
}

func TestComplexNestedSlice2(t *testing.T) {
	type Item struct {
		ItemName string `excel:"title=商品名;width=15;column=C"`
		Price    int    `excel:"title=价格;width=10;column=D"`
	}

	// User2 has 1 order with 0 items (should occupy 1 row)
	order3 := &Item{
		ItemName: "123",
		Price:    100,
	}

	_, err := XlsxWrite(nil, &order3, "DeepNested", "/Users/jeff/Desktop/deep_nested.xlsx", true)
	if err != nil {
		t.Fatalf("Failed to write deep nested struct: %v", err)
	}
}
