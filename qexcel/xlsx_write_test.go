package qexcel

import (
	"fmt"
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
	type ZXProfitExportData struct {
		RelatedMo              string `excel:"title=所属月份;width=20;column=B"`
		EntShortName           string `excel:"title=企业名称;width=20;column=C"`
		EntPayDay              string `excel:"title=发薪日;width=20;column=D"`
		EntSettleZone          string `excel:"title=发薪周期;width=20;column=E"`
		TrgtSpName             string `excel:"title=劳务名称;width=20;column=F"`
		SalaryPayer            int64  `excel:"title=月薪是否劳务打款;width=20;column=G"`
		ImportX                string `excel:"title=合计X（导入）(元);width=20;column=H"`
		EndX                   string `excel:"title=合计X（最终）(元);width=20;column=I"`
		AdjustX                string `excel:"title=X调整金额(元);width=20;column=J"`
		WeeklyPaidAmt          string `excel:"title=已发周薪(元);width=20;column=K"`
		RemainingSalary        string `excel:"title=剩余月薪(元);width=20;column=L"`
		ReturnFee              string `excel:"title=补贴金额（元）;width=20;column=M"`
		PaidTax                string `excel:"title=补贴服务费（元）;width=20;column=N"`
		ReturnFeeAfterTax      string `excel:"title=实发补贴（元）;width=20;column=O"`
		DayReturnAmount        string `excel:"title=（新）会员补贴（元）;width=20;column=P"`
		LaborY                 string `excel:"title=劳务Y（元）;width=20;column=Q"`
		SumPay                 string `excel:"title=合计支出（Y）（元）;width=20;column=R"`
		AdjustY                string `excel:"title=Y调整金额;width=20;column=R"`
		XProfit                string `excel:"title=（X-Y）（元）;width=20;column=S"`
		XProfitRatio           string `excel:"title=我打分成比例;width=20;column=T"`
		SumProfit              string `excel:"title=我打分成（元）;width=20;column=U"`
		TolWdProfit            string `excel:"title=合计我打分成（元）;width=20;column=V"`
		TolLaborProfit         string `excel:"title=合计劳务分成（元）;width=20;column=W"`
		DispatchOrderAmt       string `excel:"title=调单服务费（元）;width=20;column=X"`
		AdjustDispatchOrderAmt string `excel:"title=调单服务费调整（元）;width=20;column=Y"`
		XArrears               string `excel:"title=劳务欠款（元）;width=20;column=Z"`
		BossArrears            string `excel:"title=大佬欠款;width=20;column=AA"`
		XPaidBIllAmt           string `excel:"title=劳务到账（元）;width=20;column=AB"`
		IsCloseY               string `excel:"title=是否Y关账;width=20;column=AC"`
		IsClose                string `excel:"title=是否关账;width=20;column=AD"`
		InvoiceTrgtCn          string `excel:"title=开票单位名称;width=20;column=AE"`
		LastRunTm              string `excel:"title=最后计算时间;width=30;column=AF"`
	}

	type ZXProfitExport struct {
		BossId int64 `excel:"title=大佬;width=20;column=A"`
		Data   []*ZXProfitExportData
	}
	a := &ZXProfitExportData{
		RelatedMo:              "1",
		EntShortName:           "1",
		EntPayDay:              "1",
		EntSettleZone:          "1",
		TrgtSpName:             "",
		SalaryPayer:            0,
		ImportX:                "",
		EndX:                   "",
		AdjustX:                "",
		WeeklyPaidAmt:          "",
		RemainingSalary:        "",
		ReturnFee:              "",
		PaidTax:                "",
		ReturnFeeAfterTax:      "",
		DayReturnAmount:        "",
		LaborY:                 "",
		SumPay:                 "",
		AdjustY:                "",
		XProfit:                "",
		XProfitRatio:           "",
		SumProfit:              "",
		TolWdProfit:            "",
		TolLaborProfit:         "",
		DispatchOrderAmt:       "",
		AdjustDispatchOrderAmt: "",
		XArrears:               "",
		BossArrears:            "",
		XPaidBIllAmt:           "",
		IsCloseY:               "",
		IsClose:                "",
		InvoiceTrgtCn:          "哈哈",
		LastRunTm:              "xxxxxx",
	}
	aList := make([]*ZXProfitExportData, 0)
	aList = append(aList, a)
	aList = append(aList, a)

	b := &ZXProfitExport{
		BossId: 100,
		Data:   aList,
	}
	exportList := make([]*ZXProfitExport, 0, 10)
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
