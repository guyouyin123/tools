package qexcel

import (
	"fmt"
	"os"
	"testing"
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
		Size    int    `excel:"title=面积;width=20;column=D"`
		Address string `excel:"title=地址;width=20;column=E"`
	}
	type Like struct {
		Desc string `excel:"title=爱好;width=40;column=F"`
	}
	type Friend struct {
		Friend string `excel:"title=朋友;width=40;column=G"`
	}
	type User struct {
		Name   string `excel:"title=姓名;width=20;column=A"`
		Age    int    `excel:"title=年龄;width=40;column=B"`
		Email  string `excel:"title=身份证;width=30;column=C"`
		House  []*Info
		Like   []*Like
		Friend []*Friend
	}

	h1 := &Info{
		Size:    100,
		Address: "上海市浦东新区东方明珠",
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
