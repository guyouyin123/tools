package qtable

import (
	"io"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type execlFile struct {
	*Config
}

func newExeclFile(config *Config) *execlFile {
	return &execlFile{
		Config: config,
	}
}

func (e *execlFile) FileType() string {
	return "excel"
}

func (e *execlFile) Read(reader io.Reader) (data []map[string]string, err error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return
	}
	name := f.GetSheetName(1)
	rows, err := f.Rows(name)
	if err != nil {
		return
	}
	data = make([]map[string]string, 0)
	titles := make(map[string]int)
	initTitle := false
	rowSeq := -1
	for rows.Next() {
		rowSeq++
		columns := rows.Columns()
		if len(columns) == 0 {
			continue
		}
		// 以 # 开头的行,认为是注释行
		first := strings.TrimSpace(columns[0])
		if strings.HasPrefix(first, "#") {
			continue
		}
		// 初始化 title
		if initTitle == false {
			initTitle = true
			for k, v := range columns {
				titles[v] = k
			}
			continue
		} else {
			// 字段值小于 title,扩容 columns
			if len(columns) < len(titles) {
				expand := make([]string, len(titles)-len(columns))
				columns = append(columns, expand...)
			}
			row := make(map[string]string)
			for k := range titles {
				idx := titles[k]
				if idx >= len(titles) {
					continue
				}
				val := columns[idx]
				val = strings.Replace(val, "\t", "", -1)
				row[k] = val
			}
			data = append(data, row)
		}
	}
	return
}

type Config struct {
	// 表头是否区分大小写
	HeaderCaseSensitive bool
}
