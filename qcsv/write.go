package qcsv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"strings"
)

func ExportCsv(head []string, data [][]string) (out []byte, err error) {
	b := &bytes.Buffer{}
	f := csv.NewWriter(b)
	err = f.Write(head)
	if err != nil {
		return
	}
	err = f.WriteAll(data)
	if err != nil {
		return
	}
	return b.Bytes(), nil
}

func exportCsv(head []string, data [][]string) (out []byte, err error) {
	if len(head) == 0 {
		err = errors.New("ExportCsv Head is nil")
		return
	}
	columnCount := len(head)
	dataStr := bytes.NewBufferString("")
	//head
	for index, headElem := range head {
		separate := ","
		if index == columnCount-1 {
			separate = "\n"
		}
		dataStr.WriteString(headElem + separate)
	}
	//rows
	for _, dataArray := range data {
		if len(dataArray) != columnCount { //数据项数小于列数
			err = errors.New("数据项数小于列数")
		}
		for index, dataElem := range dataArray {
			if strings.Contains(dataElem, "\n") { // 含有换行符
				if strings.Contains(dataElem, `"`) { // 且含有"
					dataElem = strings.Replace(dataElem, `"`, `""`, -1)

				}
				dataElem = `"` + dataElem + `"`
			}
			separate := ","
			if index == columnCount-1 {
				separate = "\n"
			}
			if !strings.HasSuffix(dataElem, "\t") {
				dataElem += "\t"
			}
			dataStr.WriteString(dataElem + separate)
		}
	}
	out = make([]byte, len(dataStr.Bytes()))
	out = dataStr.Bytes()
	return
}
