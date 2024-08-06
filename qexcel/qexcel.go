package qexcel

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
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

/*
写入xlsx
不支持合并单元格
*/
func WriteToXlsxV1(dataList []interface{}, sheetName string, savePath string, isSaveFile bool) (file *xlsx.File, err error) {
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
		file.Save(savePath)
	}
	return file, nil
}

/*
写入xlsx
支持合并单元格--v2兼容v1
只支持一层嵌套
*/
func WriteToXlsxV2(dataList []interface{}, sheetName string, savePath string, isSaveFile bool) (f *excelize.File, err error) {
	//1.添加sheet
	f = excelize.NewFile()
	f.SetSheetName(sheetName, sheetName)
	if err != nil {
		return nil, err
	}

	//2.处理tag标签
	vaType := reflect.ValueOf(dataList[0]).Type()
	titleCount := vaType.NumField()
	tagMap := map[string]*Tag{}
	mergeList := make([]string, 0)
	for i := 0; i < titleCount; i++ {
		field := vaType.Field(i)
		if field.Type.Kind() == reflect.Slice {
			//子结构体
			elem := field.Type.Elem().Elem()
			elemCount := elem.NumField()
			for j := 0; j < elemCount; j++ {
				s := elem.Field(j).Tag.Get("excel")
				name := elem.Field(j).Name
				st := strings.Split(s, ";")
				if len(st) <= 1 {
					return nil, errors.New("检查结构体的excel标签")
				}
				title := strings.Split(st[0], "=")[1]
				width, _ := strconv.Atoi(strings.Split(st[1], "=")[1])
				column := strings.Split(st[2], "=")[1]
				t := &Tag{
					Title:  title,
					Width:  width,
					Column: column,
				}
				tagMap[name] = t
			}
		} else {
			s := field.Tag.Get("excel")
			st := strings.Split(s, ";")
			if len(st) <= 1 {
				return nil, errors.New("检查结构体的excel标签")
			}
			title := strings.Split(st[0], "=")[1]
			width, _ := strconv.Atoi(strings.Split(st[1], "=")[1])
			column := strings.Split(st[2], "=")[1]
			fieldName := vaType.Field(i).Name
			t := &Tag{
				Title:  title,
				Width:  width,
				Column: column,
			}
			tagMap[fieldName] = t
			mergeList = append(mergeList, column)
		}
	}

	//3.处理title
	for _, v := range tagMap {
		f.SetCellValue(sheetName, fmt.Sprintf("%s%d", v.Column, 1), v.Title)
		f.SetColWidth(sheetName, v.Column, v.Column, float64(v.Width))
	}

	//4.处理数据
	row := 2
	for i := 0; i < len(dataList); i++ {
		maxElemSliceCount := 0
		lastRow := 0
		data := dataList[i]
		va := reflect.ValueOf(data)
		vaTyp := va.Type()
		titleCount2 := va.NumField()
		isz := false
		for j := 0; j < titleCount2; j++ {
			field := vaTyp.Field(j)
			if field.Type.Kind() == reflect.Slice {
				//子结构体数组
				lastRow = row
				elemSliceObj := reflect.ValueOf(data).Field(j)
				elemSliceCount := elemSliceObj.Len()
				if elemSliceCount > maxElemSliceCount {
					maxElemSliceCount = elemSliceCount
					isz = true
				}
				for x := 0; x < elemSliceCount; x++ {
					elemObj := elemSliceObj.Index(x).Elem()
					elemObjCount := elemObj.NumField()
					// 遍历切片元素的字段
					for k := 0; k < elemObjCount; k++ {
						elemField := elemObj.Field(k)
						elemFieldType := elemObj.Type().Field(k)
						tagInfo, ok := tagMap[elemFieldType.Name]
						if !ok {
							continue
						}
						f.SetCellValue(sheetName, fmt.Sprintf("%s%d", tagInfo.Column, lastRow), elemField.Interface())
					}
					lastRow++
				}
			} else {
				//主结构
				fileName := field.Name
				vaField := va.Field(j)
				fileValue := vaField.Interface()
				tagInfo, ok := tagMap[fileName]
				if !ok {
					continue
				}
				pos := fmt.Sprintf("%s%d", tagInfo.Column, row)
				f.SetCellValue(sheetName, pos, fileValue)
			}
		}
		//preRow := row
		if isz {
			isz = false
			row = row + maxElemSliceCount
			//合并单元格
			for _, v := range mergeList {
				merSta := fmt.Sprintf("%s%d", v, row-maxElemSliceCount)
				merEnd := fmt.Sprintf("%s%d", v, row-1)
				f.MergeCell(sheetName, merSta, merEnd)
			}
		} else {
			row++
		}
	}

	//5.增加样式--左右上下居中
	style, err := f.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center"}}`)
	if err != nil {
		return nil, err
	}
	// 设置第一个sheet页的A到Z的单元格的样式
	for col := 'A'; col <= 'Z'; col++ {
		for i := 1; i <= row; i++ {
			cell := string(col) + fmt.Sprintf("%d", i)
			f.SetCellStyle(sheetName, cell, cell, style)
		}
	}

	if isSaveFile {
		if err = f.SaveAs(savePath); err != nil {
			return nil, err
		}
	}
	return f, nil
}
