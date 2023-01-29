package qcsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func Read(b []byte) ([]map[string]string, error) {
	const UTF8BOM = "\xEF\xBB\xBF"
	bf := string(b)
	if strings.HasPrefix(bf, UTF8BOM) { //截取utf-8 BOM
		bf = strings.TrimPrefix(bf, UTF8BOM)
	}
	list := make([]map[string]string, 0)
	r := csv.NewReader(strings.NewReader(bf))
	r.Comma, r.Comment = ',', '#'
	cols, err := r.Read()
	if err != nil {
		fmt.Println("Error:", err)
		return list, err
	}
	titles := make(map[string]int)
	for k, v := range cols {
		titles[v] = k
	}
	// fmt.Println("title=", titles)
	for {
		rows, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return list, err
		}
		// fmt.Println(record)
		row := make(map[string]string)
		for k := range titles {
			val := rows[titles[k]]
			val = strings.Replace(val, "\t", "", -1)
			row[k] = val
		}
		list = append(list, row)
	}
	// fmt.Println("csv成功读取数据", len(list))
	return list, nil
}
