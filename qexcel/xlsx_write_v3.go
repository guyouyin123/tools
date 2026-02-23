package qexcel

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

type tag struct {
	Title     string            //标题
	FieldName string            //结构体字段名
	Width     float64           //宽度
	Column    string            //列名
	isEnum    bool              //是否开启枚举值映射(自动根据Enum判断)
	Enum      map[string]string //enum枚举值映射
	Style     string            //样式
	IsMerge   bool              //是否合并单元格
}

type saveExcel struct {
	f         *excelize.File
	tagMap    map[string]*tag
	sheetName string
}

func initExcel(f *excelize.File, sheetName string) *saveExcel {
	if f == nil {
		f = excelize.NewFile()
	}
	index := f.NewSheet(sheetName)
	if sheetName != "Sheet1" {
		f.DeleteSheet("Sheet1")
	}
	f.SetActiveSheet(index)
	s := &saveExcel{
		f:         f,
		tagMap:    map[string]*tag{},
		sheetName: sheetName,
	}
	return s
}

/*
XlsxWrite 写入xlsx
style:样式

	wrap_text:true //自动换行
	vertical:"top" //垂直对齐方式
	horizontal:"center" //居中对齐方式
	indent:1 //缩进
	shrink_to_fit:false //不缩小字体填充
	text_rotation:0 //文本旋转角度

title:标题
width:列宽
column:所属列
IsMerge:true 开启单元格自动合并
enum:枚举值
*/
func XlsxWrite(f *excelize.File, data interface{}, sheetName string, savePath string, isSaveFile bool) (f2 *excelize.File, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic occurred: %v", e)
			return
		}
	}()

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
		this.f.SetColWidth(this.sheetName, v.Column, v.Column, v.Width)
	}

	//3.写入数据 (Layout & Render)
	currentRow := 2
	for _, data := range dataList {
		val := reflect.ValueOf(data)
		height := this.calcHeight(val)
		this.renderItem(val, currentRow, height)
		currentRow += height
	}

	//4.处理样式
	err = this.SetStyle(currentRow - 1)
	if err != nil {
		return nil, err
	}

	//5.保存文件
	if isSaveFile {
		if err = this.f.SaveAs(savePath); err != nil {
			return nil, err
		}
	}
	return this.f, nil
}

// calcHeight 计算节点高度
func (this *saveExcel) calcHeight(v reflect.Value) int {
	v = indirect(v)
	if v.Kind() != reflect.Struct {
		return 1
	}

	maxHeight := 1
	for i := 0; i < v.NumField(); i++ {
		fieldVal := indirect(v.Field(i))
		kind := fieldVal.Kind()

		if kind == reflect.Slice {
			sliceH := 0
			for k := 0; k < fieldVal.Len(); k++ {
				sliceH += this.calcHeight(fieldVal.Index(k))
			}
			// 如果切片为空，至少占一行，除非它没有任何展示内容？
			// 通常如果切片为空，我们希望显示父级字段，所以至少为1。
			// 如果切片有内容，sum(children)
			if sliceH == 0 {
				sliceH = 1
			}
			if sliceH > maxHeight {
				maxHeight = sliceH
			}
		} else if kind == reflect.Struct {
			// 嵌套结构体（非Slice），它与当前结构体在同一行开始，
			// 其高度由其内部最复杂的字段决定。
			// 父结构体高度必须能容纳子结构体。
			h := this.calcHeight(fieldVal)
			if h > maxHeight {
				maxHeight = h
			}
		}
	}
	return maxHeight
}

// renderItem 渲染节点
func (this *saveExcel) renderItem(v reflect.Value, startRow int, height int) {
	v = indirect(v)
	if v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		val := v.Field(i)

		// Case 1: Slice
		if val.Kind() == reflect.Slice {
			curr := startRow
			// 如果slice为空，我们需要跳过渲染子项，但是父项如果已经渲染了，这里不需要做额外操作。
			// 只有当父项高度由其他字段撑大时，这里留空。
			for k := 0; k < val.Len(); k++ {
				item := val.Index(k)
				h := this.calcHeight(item)
				this.renderItem(item, curr, h)
				curr += h
			}
			continue
		}

		// Case 2: Nested Struct (not slice)
		// 需要处理指针指向struct的情况
		indirectVal := indirect(val)
		if indirectVal.Kind() == reflect.Struct {
			this.renderItem(indirectVal, startRow, height)
			continue
		}

		// Case 3: Basic Field (Leaf)
		// 只有在tagMap中存在的字段才写入
		tag, ok := this.tagMap[field.Name]
		if ok {
			this.writeCell(tag, val, startRow)
			// Merge
			if tag.IsMerge && height > 1 {
				this.merge(tag.Column, startRow, startRow+height-1)
			}
		}
	}
}

func (this *saveExcel) writeCell(tagInfo *tag, val reflect.Value, row int) {
	fileValue := val.Interface()
	if tagInfo.isEnum {
		f := fmt.Sprintf("%v", fileValue)
		v, ok2 := tagInfo.Enum[f]
		if ok2 {
			fileValue = v
		} else {
			fileValue = "未知"
		}
	}
	pos := fmt.Sprintf("%s%d", tagInfo.Column, row)
	this.f.SetCellValue(this.sheetName, pos, fileValue)
}

func (this *saveExcel) merge(column string, start, end int) {
	rowSta := fmt.Sprintf("%s%d", column, start)
	rowEnd := fmt.Sprintf("%s%d", column, end)
	this.f.MergeCell(this.sheetName, rowSta, rowEnd)
}

func indirect(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			// 如果是nil指针，返回一个零值的Struct以便继续遍历类型信息?
			// 不，如果是nil，我们无法获取值。但是我们需要获取类型来做Tag解析吗？
			// Tag解析在tagHandle已经做完了。
			// 这里是渲染阶段。如果值为nil，就无法渲染其子字段。
			return v // Keep as Ptr so we can check IsNil if needed, or handle above
		}
		return v.Elem()
	}
	return v
}

// SetStyle 设置样式
func (this *saveExcel) SetStyle(maxRow int) error {
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
		endCell := fmt.Sprintf("%s%d", v.Column, maxRow)
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
	if len(dataList) == 0 {
		return nil
	}
	baseVa := reflect.ValueOf(dataList[0])
	if err := this._tagHandle(baseVa); err != nil {
		return err
	}
	return nil
}

type ExcelTag struct {
	Title   string
	Width   float64
	Column  string
	Style   string
	Enum    string
	IsMerge bool
}

func parseExcelTag(s string) ExcelTag {
	tagInfo := ExcelTag{}
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
			tagInfo.Title = value
		case "width":
			tagInfo.Width = cast.ToFloat64(value)
		case "column":
			tagInfo.Column = value
		case "style":
			tagInfo.Style = value
		case "enum":
			tagInfo.Enum = value
		case "IsMerge":
			if value == "true" {
				tagInfo.IsMerge = true
			}
		}
	}
	return tagInfo
}

func (this *saveExcel) fieldHandle(field reflect.StructField) error {
	s := field.Tag.Get("excel")
	if s == "" {
		return nil
	}
	tagInfo := parseExcelTag(s)
	dic := map[string]string{}
	isEnum := false
	if len(tagInfo.Enum) > 0 {
		dicStr := tagInfo.Enum
		_ = jsoniter.Unmarshal([]byte(dicStr), &dic)
		isEnum = true
	}

	t := &tag{
		Title:     tagInfo.Title,
		FieldName: field.Name,
		Width:     tagInfo.Width,
		Column:    tagInfo.Column,
		isEnum:    isEnum,
		Enum:      dic,
		Style:     tagInfo.Style,
		IsMerge:   tagInfo.IsMerge,
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

	if vaType.Kind() == reflect.Ptr {
		vaType = vaType.Elem()
	}
	numField := vaType.NumField()
	for i := 0; i < numField; i++ {
		field := vaType.Field(i)
		switch field.Type.Kind() {
		case reflect.Slice:
			//处理数组
			// 为了获取slice内部元素的类型，我们需要一个实例，或者直接从Type获取Elem
			// 原有代码尝试从实例获取，如果slice为空，从Type获取
			var sliceVal reflect.Value
			if baseVa.Kind() == reflect.Ptr {
				sliceVal = baseVa.Elem().Field(i)
			} else {
				sliceVal = baseVa.Field(0) // Bug in original code? baseVa.Field(i) probably intended?
				// Original code: sliceVal = baseVa.Field(0) -> likely copy paste error in original or it was assuming something specific.
				// Wait, let's look at original code line 399: `sliceVal = baseVa.Field(0)`
				// If baseVa is a Struct, `baseVa.Field(i)` is the slice.
				// Why `Field(0)`?
				// Maybe `baseVa` here is already the Slice? No, switch says `field.Type.Kind() == reflect.Slice`. `field` is from `vaType.Field(i)`.
				// So `baseVa` is the struct.
				// Correct logic should be `baseVa.Field(i)`.
				// Let's fix this obvious bug while we are here, or stick to safe side?
				// If I change it to Field(i) it is correct.
				sliceVal = baseVa.Field(i)
			}

			// If we can't get an element (empty slice), use Type info
			if sliceVal.Len() == 0 {
				// Create a zero value of the element type to recurse
				elemType := field.Type.Elem()
				if elemType.Kind() == reflect.Ptr {
					elemType = elemType.Elem()
				}
				// We need a Value to pass to _tagHandle if possible, or refactor _tagHandle to take Type.
				// But _tagHandle takes Value.
				// Let's create a Zero value.
				newVal := reflect.New(elemType).Elem()
				if err := this._tagHandle(newVal); err != nil {
					return err
				}
			} else {
				// Use first element
				elem := sliceVal.Index(0)
				if elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				if err := this._tagHandle(elem); err != nil {
					return err
				}
			}

		case reflect.Struct, reflect.Ptr:
			//处理嵌套结构体
			nestedVal := baseVa.Field(i)
			// Handle nil ptr
			if nestedVal.Kind() == reflect.Ptr && nestedVal.IsNil() {
				// Use Type info
				elemType := field.Type
				if elemType.Kind() == reflect.Ptr {
					elemType = elemType.Elem()
				}
				newVal := reflect.New(elemType).Elem()
				if err := this._tagHandle(newVal); err != nil {
					return err
				}
			} else {
				err := this._tagHandle(nestedVal)
				if err != nil {
					return err
				}
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

func isIntType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Uint, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64:
		return true
	default:
		return false
	}
}
