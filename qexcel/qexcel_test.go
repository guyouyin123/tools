package qexcel

import (
	"fmt"
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

	file, err := WriteToXlsxV1(list, "Sheet1", "./test.xlsx", true)
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

	file, err := WriteToXlsxV2(list, "Sheet1", "./test1.xlsx", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(file)
}

func TestDemo(t *testing.T) {
	type NameListExportResp struct {
		NameListId int64  `excel:"title=会员ID;width=20;column=A"` //会员ID
		IntvDt     string `excel:"title=面试时间;width=20;column=B"` //面试时间
		EntryDt    string `excel:"title=入职时间;width=20;column=C"` //入职时间
		LeaveDt    string `excel:"title=离职时间;width=20;column=D"` //离职时间
		IntvSts    string `excel:"title=面试状态;width=20;column=E"` //面试状态 0 未处理， 1未面试 2 面试通过 3 面试不通过 4 放弃
		WorkSts    string `excel:"title=工作状态;width=20;column=F"` //工作状态，1 在职 2 离职 3 转正 4 未处理 5 未知 6 自离
	}

	n := &NameListExportResp{
		NameListId: 123,
		IntvDt:     "2024-01-01",
		EntryDt:    "2024-01-01",
		LeaveDt:    "0001-01-01",
		IntvSts:    "未处理",
		WorkSts:    "在职",
	}
	list := []interface{}{}
	list = append(list, n)
	file, err := WriteToXlsxV2(list, "Sheet1", "./test2.xlsx", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(file)
}
