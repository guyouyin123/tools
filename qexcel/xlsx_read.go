package qexcel

import (
	"github.com/tealeg/xlsx"
)

// 不支持行合并的excel
func XlsxRead(tablename string) ([]map[string]interface{}, error) {
	xlFile, err := xlsx.OpenFile(tablename) //所有的sheet
	if err != nil {
		return nil, err
	}
	titleMap := map[int]string{}
	titleRow := xlFile.Sheets[0].Rows[0].Cells
	for k, titlename := range titleRow {
		titleMap[k] = titlename.String()
	}
	T := []map[string]interface{}{}
	for _, sheet := range xlFile.Sheets {
		//sheet表
		for i := 1; i < len(sheet.Rows); i++ {
			//表中的行
			t := map[string]interface{}{}
			for k, v := range sheet.Rows[i].Cells {
				//行中的列
				t[titleMap[k]] = v.Value
			}
			T = append(T, t)
		}
	}
	return T, nil
}
