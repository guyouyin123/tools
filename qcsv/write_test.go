package qcsv

import (
	"fmt"
	"os"
	"testing"
)

func TestExportCsv(t *testing.T) {
	var csvD [][]string
	header := []string{"姓名", "身份证号码", "手机号码", "地址"}
	csvD = append(csvD, []string{
		"测试123,123,123", "12312312312", "13223213123123", `123asda`, "sfsdf",
	})
	data, err := ExportCsv(header, csvD)
	if err != nil {
		t.Fatal(err)
	}
	temp, err := os.Create("temp.csv") //TODO:高并发时不安全
	if err != nil {
		fmt.Println(err)
	}
	defer temp.Close()
	temp.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	temp.Write(data)
}
