package qexcel

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"strconv"
	"strings"
)

/*
XlsxWriteV3 写入xlsx
支持合并单元格--v3兼容v1,V2
v2只支持1层嵌套，只支持excel一层合并。v3支持无级嵌套，无级合并
结构体类型支持指针和非指针
*/

type Tag struct {
	Title     string            //标题
	FieldName string            //结构体字段名
	Width     int               //宽度
	Column    string            //列名
	isEnum    bool              //是否开启枚举值映射(自动根据Enum判断)
	Enum      map[string]string //enum枚举值映射
	Style     string            //样式
}

type saveExcel struct {
	f         *excelize.File
	tagMap    map[string]*Tag
	mergeMap  map[string][][2]int //合并单元格map{"A":[[1,4],[7,10]]}
	sheetName string
	row       int
	addRow    bool
}

func initExcel(f *excelize.File, sheetName string) *saveExcel {
	if f == nil {
		f = excelize.NewFile()
	}
	index := f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)
	s := &saveExcel{
		f:         f,
		tagMap:    map[string]*Tag{},
		mergeMap:  map[string][][2]int{},
		sheetName: sheetName,
		row:       2,
		addRow:    false,
	}
	return s
}

/*
XlsxWriteV3 写入xlsx
ps:数组结构体需要放在最后
错误事例：

	type User struct {
		Name  string `excel:"title=姓名;width=20;column=F"`
		Class []*Class
		Age   int `excel:"title=年龄;width=20;column=B"`
	}

正确事例：

	type like struct {
			Like string `excel:"title=爱好;width=20;column=C;style={\"alignment\":{\"horizontal\":\"center\",\"vertical\":\"center\",\"wrap_text\":true}}"`
		}
		type User struct {
			Name          string `excel:"title=姓名;width=20;column=A;style={\"alignment\":{\"horizontal\":\"center\",\"text_rotation\":45}}"`
			Age           int    `excel:"title=年龄;width=20;column=B;"`
			SettlementTyp int8   `excel:"title=模式;width=10;column=F;enum={\"0\":\"未知\",\"1\":\"ZX模式\",\"2\":\"Z模式\",\"3\":\"ZA模式\",\"4\":\"Z-B模式\",\"5\":\"ZX-B模式\",\"6\":\"ZX-A模式\",\"7\":\"Z-D模式\",\"8\":\"ZX-D模式\"}"`
			Like          []*like
		}

style:

	wrap_text:true //自动换行
	vertical:"top" //垂直对齐方式
	horizontal:"center" //居中对齐方式
	indent:1 //缩进
	shrink_to_fit:false //不缩小字体填充
	text_rotation:0 //文本旋转角度
*/
func XlsxWriteV3(f *excelize.File, data interface{}, sheetName string, savePath string, isSaveFile bool) (f2 *excelize.File, err error) {
	dataList := make([]interface{}, 0)
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v := reflect.ValueOf(data).Elem()
		if v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				dataList = append(dataList, v.Index(i).Interface())
			}
		}
	} else if t.Kind() == reflect.Slice {
		v := reflect.ValueOf(data)
		for i := 0; i < v.Len(); i++ {
			dataList = append(dataList, v.Index(i).Interface())
		}
	} else {
		dataList = append(dataList, data)
	}

	this := initExcel(f, sheetName)
	if len(dataList) == 0 {
		return f, nil
	}

	//1.处理tag标签
	err = this.tagHandle(dataList)
	if err != nil {
		return nil, err
	}
	//2.处理title
	for _, v := range this.tagMap {
		this.f.SetCellValue(this.sheetName, fmt.Sprintf("%s%d", v.Column, 1), v.Title)
		this.f.SetColWidth(this.sheetName, v.Column, v.Column, float64(v.Width))
	}

	//3.写入数据
	this.WriteDate(dataList)

	//4.处理合并单元格
	this.MergeCell()
	//5.处理样式
	err = this.SetStyle()
	if err != nil {
		return nil, err
	}

	//6.保存文件
	if isSaveFile {
		if err = this.f.SaveAs(savePath); err != nil {
			return nil, err
		}
	}
	return this.f, nil
}

// MergeCell 合并单元格
func (this *saveExcel) MergeCell() {
	//双指针
	sta := 0
	end := 0
	for _, tag := range this.tagMap {
		for i := 2; i <= this.row; i++ {
			cell := fmt.Sprintf("%s%d", tag.Column, i)
			value := this.f.GetCellValue(this.sheetName, cell)
			if value != "" {
				if end != 0 && end > sta {
					this.mergeMap[tag.Column] = append(this.mergeMap[tag.Column], [2]int{sta, end})
					sta, end = 0, 0
				}
				sta = i
			} else {
				end = i
				if end > sta && i == this.row {
					this.mergeMap[tag.Column] = append(this.mergeMap[tag.Column], [2]int{sta, end})
					sta, end = 0, 0
				}
			}
		}
	}
	for column, mergeL := range this.mergeMap {
		for _, merge := range mergeL {
			this.f.MergeCell(this.sheetName, fmt.Sprintf("%s%d", column, merge[0]), fmt.Sprintf("%s%d", column, merge[1]))
		}
	}
}

// SetStyle 设置样式
func (this *saveExcel) SetStyle() error {
	// 设置单元格的样式
	//设置列样式
	for _, v := range this.tagMap {
		styleStr := this.tagMap[v.FieldName].Style
		if len(styleStr) == 0 {
			continue
		}
		style, err := this.f.NewStyle(styleStr)
		if err != nil {
			return err
		}
		startCell := fmt.Sprintf("%s1", v.Column)
		endCell := fmt.Sprintf("%s%d", v.Column, this.row)
		this.f.SetCellStyle(this.sheetName, startCell, endCell, style)
	}

	//设置标题加粗
	style2, err := this.f.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center"},"font": {"bold": true}}`)
	if err != nil {
		return err
	}
	for _, v := range this.tagMap {
		title := fmt.Sprintf("%s1", v.Column)
		this.f.SetCellStyle(this.sheetName, title, title, style2)
	}
	return nil
}

// tagHandle 标签处理
func (this *saveExcel) tagHandle(dataList []interface{}) error {
	baseVa := reflect.ValueOf(dataList[0])
	if err := this._tagHandle(baseVa); err != nil {
		return err
	}
	return nil
}

type ExcelTag struct {
	Title  string
	Width  int
	Column string
	Style  string
	Enum   string
}

func parseExcelTag(s string) ExcelTag {
	tag := ExcelTag{}
	pairs := strings.Split(s, ";")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) < 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "title":
			tag.Title = value
		case "width":
			if width, err := strconv.Atoi(value); err == nil {
				tag.Width = width
			}
		case "column":
			tag.Column = value
		case "style":
			tag.Style = value
		case "enum":
			tag.Enum = value
		}
	}
	return tag
}

func (this *saveExcel) fieldHandle(field reflect.StructField) error {
	s := field.Tag.Get("excel")
	tag := parseExcelTag(s)
	dic := map[string]string{}
	isEnum := false
	if len(tag.Enum) > 0 {
		dicStr := tag.Enum
		_ = jsoniter.Unmarshal([]byte(dicStr), &dic)
		isEnum = true
	}

	t := &Tag{
		Title:     tag.Title,
		FieldName: field.Name,
		Width:     tag.Width,
		Column:    tag.Column,
		isEnum:    isEnum,
		Enum:      dic,
		Style:     tag.Style,
	}
	this.tagMap[field.Name] = t
	return nil
}

func (this *saveExcel) _tagHandle(baseVa reflect.Value) error {
	var vaType reflect.Type
	if baseVa.Kind() == reflect.Ptr {
		vaType = baseVa.Type().Elem()
	} else {
		vaType = baseVa.Type()
	}
start:
	if vaType.Kind() == reflect.Ptr {
		vaType = vaType.Elem()
	}
	numField := vaType.NumField()
	for i := 0; i < numField; i++ {
		field := vaType.Field(i)
		switch field.Type.Kind() {
		case reflect.Slice:
			//处理数组
			var sliceVal reflect.Value
			if baseVa.Kind() == reflect.Ptr {
				sliceVal = baseVa.Elem().Field(i)
			} else {
				sliceVal = baseVa.Field(i)
			}
			le := sliceVal.Len()
			if le == 0 {
				vaType = field.Type.Elem()
				goto start
			} else {
				baseVa = sliceVal.Index(0)
				if baseVa.Kind() == reflect.Ptr {
					baseVa = baseVa.Elem()
				}
				err := this._tagHandle(baseVa)
				if err != nil {
					return err
				}
			}
		case reflect.Struct, reflect.Ptr:
			//处理嵌套结构体
			nestedVal := baseVa.Field(i)
			err := this._tagHandle(nestedVal)
			if err != nil {
				return err
			}
		default:
			err := this.fieldHandle(field)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *saveExcel) WriteDate(dataList []interface{}) {
	for index, data := range dataList {
		var va reflect.Value
		baseDataVa := reflect.ValueOf(data)
		if baseDataVa.Kind() == reflect.Ptr {
			va = baseDataVa.Elem()
		} else {
			va = baseDataVa
		}
		vaTyp := va.Type()
		titleCount := va.NumField()
		for i := 0; i < titleCount; i++ {
			field := vaTyp.Field(i)
			vaField := va.Field(i)
			elemSliceObj := reflect.Value{}
			if baseDataVa.Kind() == reflect.Ptr {
				//子结构为指针
				elemSliceObj = baseDataVa.Elem().Field(i)
			} else {
				//子结构为非指针
				elemSliceObj = baseDataVa.Field(i)
			}
			this.writeExcel(field, elemSliceObj, vaField)
		}
		if index != len(dataList)-1 {
			this.row++
		}
	}
}

func (this *saveExcel) writeExcel(field reflect.StructField, elemSliceObj, vaField reflect.Value) {
	if field.Type.Kind() == reflect.Slice {
		// 子结构体数组
		elemSliceCount := 0
		if elemSliceObj.Kind() != reflect.Slice {
			fileValue := vaField.Interface()
			arrayValue := reflect.ValueOf(fileValue)
			elemSliceCount = arrayValue.Len()
			elemSliceObj = arrayValue
			this.addRow = false
		} else {
			elemSliceCount = elemSliceObj.Len()
		}
		for x := 0; x < elemSliceCount; x++ {
			elemObj := elemSliceObj.Index(x)
			if elemObj.Kind() == reflect.Ptr {
				elemObj = elemObj.Elem()
			}
			for i := 0; i < elemObj.NumField(); i++ {
				field2 := elemObj.Type().Field(i)
				vaField2 := elemObj.Field(i)
				this.writeExcel(field2, elemObj, vaField2)
			}
			this.addRow = true
		}
		this.addRow = false
	} else {
		this.write(field, vaField)
	}
}

func (this *saveExcel) write(field reflect.StructField, vaField reflect.Value) {
	fileValue := vaField.Interface()
	fieldName := field.Name
	tagInfo, ok := this.tagMap[fieldName]
	if !ok {
		return
	}
	if this.addRow {
		this.row++
	}
	if tagInfo.isEnum {
		f := fmt.Sprintf("%v", fileValue)
		v, ok2 := tagInfo.Enum[f]
		if ok2 {
			fileValue = v
		} else {
			fileValue = "未知"
		}
	}
	pos := fmt.Sprintf("%s%d", tagInfo.Column, this.row)
	this.f.SetCellValue(this.sheetName, pos, fileValue)
}

func isIntType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
		return true
	default:
		return false
	}
}
