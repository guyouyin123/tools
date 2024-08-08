package qexcel

import (
	"bytes"
	"encoding/csv"
)

func CsvWrite(head []string, data [][]string) (out []byte, err error) {
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
