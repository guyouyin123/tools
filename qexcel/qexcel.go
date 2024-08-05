package qexcel

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"reflect"
	"strconv"
	"strings"
)

type Tag struct {
	Title  string
	Width  int
	Column string
}

func WriteToXlsx(dataList []interface{}, sheetName string, isSaveFile bool) (file *xlsx.File, err error) {
	//1.添加sheet
	file = xlsx.NewFile()
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		return nil, err
	}
	//2.处理tag标签
	t := reflect.TypeOf(dataList[0])
	titleCount := t.NumField()
	tagMap := map[string]*Tag{}
	for i := 0; i < titleCount; i++ {
		s := t.Field(i).Tag.Get("excel")
		st := strings.Split(s, ";")
		if len(st) <= 1 {
			return nil, errors.New("检查结构体的excel标签")
		}
		title := strings.Split(st[0], "=")[1]
		width, _ := strconv.Atoi(strings.Split(st[1], "=")[1])
		column := strings.Split(st[2], "=")[1]
		tag := &Tag{
			Title:  title,
			Width:  width,
			Column: column,
		}
		tagMap[column] = tag
	}

	titleColumnList := []string{}
	maxAscii := 'A'
	for column := range tagMap {
		columnAscii := []rune(column)[0]
		if columnAscii > maxAscii && columnAscii <= 'Z' {
			maxAscii = columnAscii
		}
	}
	for i := 'A'; i <= maxAscii; i++ {
		titleColumnList = append(titleColumnList, fmt.Sprintf("%c", i))
	}

	nilMap := map[int]struct{}{}
	//3.处理title
	row := sheet.AddRow()
	cell := row.AddCell()
	for index, titleTypeName := range titleColumnList {
		tag, ok := tagMap[titleTypeName]
		if !ok {
			nilMap[index] = struct{}{}
		} else {
			sheet.SetColWidth(index, index, float64(tag.Width))
			cell.Value = tag.Title
		}
		cell = row.AddCell()
	}

	//4.处理数据
	for i := 0; i < len(dataList); i++ {
		r := dataList[i]
		nextRow := sheet.AddRow()
		s := reflect.ValueOf(r)
		c := 0
		for j := 0; j < len(titleColumnList); j++ {
			_, ok := nilMap[j]
			if ok {
				cell = nextRow.AddCell()
			} else {
				cell = nextRow.AddCell()
				value := s.Field(c)
				cell.Value = fmt.Sprintf("%v", value)
				c++
			}
		}
	}
	if isSaveFile {
		file.Save("./test.xlsx")
	}
	return file, nil
}

type User struct {
	Name string `excel:"title=姓名;width=20;column=A"`
	Age  int    `excel:"title=年龄;width=50;column=B"`
}
