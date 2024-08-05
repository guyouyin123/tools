package qexcel

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tealeg/xlsx"
	excelize2 "github.com/xuri/excelize/v2"
)

func WriteXlsx1() {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	sheet.SetColWidth(0, 1, 20)
	cell := sheet.Cell(1, 1)
	cell.SetValue("Hello, world!")
	file.Save("./aaa.xlsx")
}

// f := excelize.NewFile()
func WriteXlsx2() {
	f := excelize2.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.SetCellValue("Sheet1", "B2", 100)
	f.SetActiveSheet(index)
	if err := f.SaveAs("./bbb.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func WriteXlsx3() {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("人员信息收集")
	if err != nil {
		panic(err.Error())
	}
	row := sheet.AddRow()
	cell := row.AddCell()
	cell.Value = "姓名"
	cell = row.AddCell()
	cell.Value = "性别"

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "张三"
	cell = row.AddCell()
	cell.Value = "男"

	err = file.Save("./ccc.xlsx")
	if err != nil {
		panic(err.Error())
	}
}

func WriteXlsx4() {
	jeff := User{
		Name: "jeff",
		Age:  18,
	}
	chary := User{
		Name: "chary",
		Age:  20,
	}
	list := []User{}
	list = append(list, jeff, chary)

	// 创建 Excel 文件
	f := excelize.NewFile()

	// 设置工作表名称
	f.SetSheetName("Sheet1", "用户信息")

	// 写入表头
	f.SetCellValue("用户信息", "A1", "姓名")
	f.SetCellValue("用户信息", "B1", "年龄")

	// 写入数据
	row := 2
	for _, user := range list {
		f.SetCellValue("用户信息", fmt.Sprintf("A%d", row), user.Name)
		f.SetCellValue("用户信息", fmt.Sprintf("B%d", row), user.Age)
		row++
	}

	// 保存 Excel 文件
	if err := f.SaveAs("./ddd.xlsx"); err != nil {
		fmt.Println(err)
		return
	}

	return
}
