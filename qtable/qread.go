package qtable

import (
	"fmt"
	"io"
)

// 读取表格文件

type TableReader interface {
	FileType() string
	Read(reader io.Reader) ([]map[string]string, error)
	// Header() []string
	// RawHeader() []Column
}

func TableRead(ext string) (fr TableReader, err error) {
	switch ext {
	case ".csv":
		return newCsvRead(), nil
	case ".xlsx":
		return newExeclFile(&Config{}), nil
	}
	return nil, fmt.Errorf("table reader %s not exist", ext)
}
