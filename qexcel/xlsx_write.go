package qexcel

// 本文件是 XlsxWrite(见 xlsx_write_v2.go)的共享底座: tag 解析、布局高度计算、反射解引用等,
// 均与具体 excelize 版本无关。旧的 XlsxWrite(基于 excelize v1)已删除, 统一用 XlsxWrite 替代。

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

// cellRef 拼单元格坐标(如 "A"+2 -> "A2")。
// 替代逐格 fmt.Sprintf("%s%d",...): 几十万行×多列会调用上百万次, Sprintf 的反射+格式化开销可观。
func cellRef(col string, row int) string {
	return col + strconv.Itoa(row)
}

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

// fieldKey 唯一标识一个字段: 用「所属结构体类型 + 字段名」联合做 key,
// 避免不同嵌套结构体里的同名字段(如两个结构体都有 Status)在 tagMap 里互相覆盖。
type fieldKey struct {
	owner reflect.Type
	name  string
}

// saveExcel 承载一次导出的字段 tag 映射(由 tagHandle 填充), 供布局高度计算复用。
type saveExcel struct {
	tagMap map[fieldKey]*tag
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

// indirect 循环解引用到非指针为止(支持 *T、**T 等任意层级指针)。
// 只解一层的话, 传入 **Struct 时会停在 *Struct(非 struct), 导致 renderItem 直接跳过、只出标题无内容。
// 遇到 nil 指针则原样返回该指针(渲染阶段无法再取其子字段)。
func indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}
	return v
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

// defaultColWidth tag 未写 width 时的默认列宽(避免 SetColWidth(0) 把整列压成隐藏)
const defaultColWidth = 20.0

func parseExcelTag(s string) ExcelTag {
	tagInfo := ExcelTag{}
	hasWidth := false
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
			hasWidth = true
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
	if !hasWidth {
		tagInfo.Width = defaultColWidth
	}
	return tagInfo
}

func (this *saveExcel) fieldHandle(owner reflect.Type, field reflect.StructField) error {
	s := field.Tag.Get("excel")
	if s == "" {
		return nil
	}
	tagInfo := parseExcelTag(s)
	// column 是写入定位的必要信息, 缺失会导致非法单元格引用而静默写不进去, 这里直接报错。
	if tagInfo.Column == "" {
		return fmt.Errorf("字段 %s.%s 的 excel tag 缺少 column", owner.Name(), field.Name)
	}
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
	this.tagMap[fieldKey{owner: owner, name: field.Name}] = t
	return nil
}

func (this *saveExcel) _tagHandle(baseVa reflect.Value) error {
	// 先把 baseVa 完整解引用到具体结构体值(支持 *T / **T ... 任意层级指针),
	// nil 指针则用元素类型的零值实例代替, 以便仍能按"类型"解析出 tag。
	// 这样下面各分支统一用 baseVa.Field(i), 不必再逐处处理指针(否则多级指针会 panic)。
	baseVa = derefToStruct(baseVa)
	if baseVa.Kind() != reflect.Struct {
		return nil
	}

	vaType := baseVa.Type()
	numField := vaType.NumField()
	for i := 0; i < numField; i++ {
		field := vaType.Field(i)
		switch field.Type.Kind() {
		case reflect.Slice:
			// 切片视为"子记录集合": 有元素就取第 0 个, 空切片则用元素类型零值实例, 都递归解析。
			sliceVal := baseVa.Field(i)
			var elem reflect.Value
			if sliceVal.Len() > 0 {
				elem = sliceVal.Index(0)
			} else {
				elem = reflect.New(field.Type.Elem()).Elem()
			}
			if err := this._tagHandle(elem); err != nil {
				return err
			}
		case reflect.Struct, reflect.Ptr:
			//处理嵌套结构体(nil 指针在递归入口用零值实例兜底)
			if err := this._tagHandle(baseVa.Field(i)); err != nil {
				return err
			}
		default:
			if err := this.fieldHandle(vaType, field); err != nil {
				return err
			}
		}
	}
	return nil
}

// derefToStruct 把值解引用到非指针为止; 途中遇到 nil 指针, 用其元素类型的零值实例继续,
// 保证最终能拿到一个可 .Field() 的具体值(用于按类型解析 tag)。
func derefToStruct(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v = reflect.New(v.Type().Elem()).Elem()
			continue
		}
		v = v.Elem()
	}
	return v
}

// XlsxWrite 是 XlsxWrite 的高性能版: 基于 excelize/v2 的 StreamWriter 顺序流式写入,
// 面向几十万行大数据 —— 边写边刷、不在内存里建整棵单元格树, 耗时与内存都远低于 v1 的
// 全内存 SetCellValue 模型。入参与 XlsxWrite 对齐: 首个 f *excelize.File 可传入(为 nil 则新建),
// 用于往已有工作簿里追加一个新 sheet。
//
// 与 v1 (XlsxWrite) 的差异:
//   - 依赖 github.com/xuri/excelize/v2, 传入/返回的 *excelize.File 都是 v2 类型(与 v1 的不通用)。
//   - tag 解析 / 布局高度 / 合并单元格 / 枚举映射 与 v1 完全一致(复用同一套逻辑);
//     tag 里的 style JSON 串也兼容(内部做 snake_case -> excelize.Style 的显式转换)。
//
// StreamWriter 固有约束: 行必须自上而下、每行恰好写一次, 且目标 sheet 在流式写入前应为空。
// 因为本包的嵌套布局里父级单元格要跨子行合并(不是天然按行产生), 所以内部先把整表渲染进"行矩阵",
// 再按行号递增顺序 SetRow 吐出。若传入已有文件, 请确保 sheetName 是一个新的/空的 sheet。
func XlsxWrite(f *excelize.File, data interface{}, sheetName, savePath string, isSaveFile bool) (f2 *excelize.File, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic occurred: %v", e)
		}
	}()

	// 1. 摊平入参为一维数据列表(指针/切片/单值均可, 与 v1 同义)
	dataList := flattenData(data)
	if len(dataList) == 0 {
		if f, err = prepareSheet(f, sheetName); err != nil {
			return nil, err
		}
		if isSaveFile {
			if err = f.SaveAs(savePath); err != nil {
				return nil, err
			}
		}
		return f, nil
	}

	// 2. 解析 tag(复用 v1 的 saveExcel.tagHandle, 它只用反射、不碰 excelize 文件)
	//    放在动文件之前: tag 出错(如缺 column)时直接返回, 不污染传入的 *File。
	s := &saveExcel{tagMap: map[fieldKey]*tag{}}
	if err = s.tagHandle(dataList); err != nil {
		return nil, err
	}
	if len(s.tagMap) == 0 {
		return nil, fmt.Errorf("没有可导出的字段: 结构体上没有任何 excel tag")
	}

	// 3. 列字母 -> 1based 序号, 并求最大列(行矩阵按最大列定宽)
	colIdx := make(map[string]int, len(s.tagMap))
	maxCol := 0
	for _, tg := range s.tagMap {
		ci, e := excelize.ColumnNameToNumber(tg.Column)
		if e != nil {
			return nil, fmt.Errorf("非法列名 %q: %w", tg.Column, e)
		}
		colIdx[tg.Column] = ci
		if ci > maxCol {
			maxCol = ci
		}
	}

	// 4. 先算每条数据高度(用于给行矩阵定尺寸), 再把值渲染进矩阵
	heights := make([]int, len(dataList))
	totalRows := 0
	for i, d := range dataList {
		h := s.calcHeight(reflect.ValueOf(d))
		heights[i] = h
		totalRows += h
	}
	gridRows := 1 + totalRows // 第 1 行为表头
	g := &gridWriter{tagMap: s.tagMap, s: s, colIdx: colIdx, grid: make([][]interface{}, gridRows)}
	for i := range g.grid {
		g.grid[i] = make([]interface{}, maxCol)
	}
	for _, tg := range s.tagMap { // 表头
		g.grid[0][colIdx[tg.Column]-1] = tg.Title
	}
	cur := 2 // 数据从第 2 行起
	for i, d := range dataList {
		g.render(reflect.ValueOf(d), cur, -1)
		cur += heights[i]
	}

	// 5. 准备文件与目标 sheet(可传入已有 *File 追加; nil 则新建), 再流式写出
	if f, err = prepareSheet(f, sheetName); err != nil {
		return nil, err
	}
	sw, e := f.NewStreamWriter(sheetName)
	if e != nil {
		return nil, e
	}

	// 5a. 列宽 + 预建样式(SetColWidth 必须在 SetRow 之前; 样式只 NewStyle 一次后复用)
	headerStyle, e := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Bold: true},
	})
	if e != nil {
		return nil, e
	}
	colStyle := make([]int, maxCol+1) // 1based 列 -> styleID(0 表示无列样式)
	for _, tg := range s.tagMap {
		ci := colIdx[tg.Column]
		if e = sw.SetColWidth(ci, ci, tg.Width); e != nil {
			return nil, e
		}
		if tg.Style != "" {
			st, e2 := parseStyleJSON(tg.Style)
			if e2 != nil {
				return nil, fmt.Errorf("列 %s 样式解析失败: %w", tg.Column, e2)
			}
			id, e3 := f.NewStyle(st)
			if e3 != nil {
				return nil, e3
			}
			colStyle[ci] = id
		}
	}

	// 5b. 逐行 SetRow(StreamWriter 要求行号递增、一次一整行; 用 Cell 携带样式)
	for r := 1; r <= gridRows; r++ {
		rowVals := make([]interface{}, maxCol)
		for c := 0; c < maxCol; c++ {
			style := colStyle[c+1]
			if r == 1 {
				style = headerStyle // 表头用加粗居中样式
			}
			rowVals[c] = excelize.Cell{StyleID: style, Value: g.grid[r-1][c]}
		}
		if e = sw.SetRow(cellRef("A", r), rowVals); e != nil {
			return nil, e
		}
	}

	// 5c. 合并单元格(值已落在块首行, 被合并掉的下方单元格为空)
	for _, m := range g.merges {
		if e = sw.MergeCell(cellRef(m.col, m.start), cellRef(m.col, m.end)); e != nil {
			return nil, e
		}
	}

	if e = sw.Flush(); e != nil {
		return nil, e
	}

	// 6. 保存
	if isSaveFile {
		if err = f.SaveAs(savePath); err != nil {
			return nil, err
		}
	}
	return f, nil
}

// prepareSheet 准备输出文件与目标 sheet: f 为 nil 则新建; 确保 sheetName 存在并设为激活 sheet;
// 仅当文件是自己新建的、且目标非默认 Sheet1 时, 删掉 NewFile 自带的空 Sheet1(传入的已有文件不动其 Sheet1)。
func prepareSheet(f *excelize.File, sheetName string) (*excelize.File, error) {
	created := false
	if f == nil {
		f = excelize.NewFile()
		created = true
	}
	if _, err := f.NewSheet(sheetName); err != nil { // 已存在则返回其索引, 不存在则新建
		return nil, err
	}
	if created && sheetName != "Sheet1" {
		if err := f.DeleteSheet("Sheet1"); err != nil {
			return nil, err
		}
	}
	if idx, err := f.GetSheetIndex(sheetName); err == nil { // 删表后索引会变, 重新取一次
		f.SetActiveSheet(idx)
	}
	return f, nil
}

// flattenData 把入参摊平成一维数据列表: *[]T / []T / *T / T 皆可(与 v1 XlsxWrite 的分支同义)。
func flattenData(data interface{}) []interface{} {
	t := reflect.TypeOf(data)
	if t == nil {
		return nil
	}
	switch t.Kind() {
	case reflect.Ptr:
		v := reflect.ValueOf(data).Elem()
		if v.Kind() == reflect.Slice {
			out := make([]interface{}, 0, v.Len())
			for i := 0; i < v.Len(); i++ {
				out = append(out, v.Index(i).Interface())
			}
			return out
		}
		return []interface{}{data}
	case reflect.Slice:
		v := reflect.ValueOf(data)
		out := make([]interface{}, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			out = append(out, v.Index(i).Interface())
		}
		return out
	default:
		return []interface{}{data}
	}
}

// gridMerge 记录一次待合并区间(同列的 start..end 行)。
type gridMerge struct {
	col        string
	start, end int
}

// gridWriter 把嵌套结构体渲染进行矩阵 grid, 逻辑与 v1 的 renderItem 一致, 只是不直接写文件。
type gridWriter struct {
	tagMap map[fieldKey]*tag
	s      *saveExcel
	colIdx map[string]int
	grid   [][]interface{} // grid[r-1][c-1] = excel 第 r 行第 c 列的值
	merges []gridMerge
}

// render 从 startRow 起渲染 v, 返回本节点占用的固有高度。语义同 xlsx_write.go 的 renderItem。
func (g *gridWriter) render(v reflect.Value, startRow, mergeHeight int) int {
	v = indirect(v)
	if v.Kind() != reflect.Struct {
		return 1
	}
	t := v.Type()
	n := v.NumField()

	// 第一遍: 渲染切片子记录(递归返回高度累加), 算出本节点固有高度
	height := 1
	for i := 0; i < n; i++ {
		val := v.Field(i)
		if val.Kind() == reflect.Slice {
			curr := startRow
			for k := 0; k < val.Len(); k++ {
				curr += g.render(val.Index(k), curr, -1)
			}
			if span := curr - startRow; span > height {
				height = span
			}
			continue
		}
		if iv := indirect(val); iv.Kind() == reflect.Struct {
			if h := g.s.calcHeight(iv); h > height {
				height = h
			}
		}
	}
	if mergeHeight <= 0 {
		mergeHeight = height
	}

	// 第二遍: 写叶子(按 mergeHeight 记录合并)与非切片嵌套结构体(共用父块高)
	for i := 0; i < n; i++ {
		field := t.Field(i)
		val := v.Field(i)
		if val.Kind() == reflect.Slice {
			continue
		}
		if iv := indirect(val); iv.Kind() == reflect.Struct {
			g.render(iv, startRow, mergeHeight)
			continue
		}
		if tg, ok := g.tagMap[fieldKey{owner: t, name: field.Name}]; ok {
			g.setCell(tg, val, startRow)
			if tg.IsMerge && mergeHeight > 1 {
				g.merges = append(g.merges, gridMerge{tg.Column, startRow, startRow + mergeHeight - 1})
			}
		}
	}
	return height
}

func (g *gridWriter) setCell(tg *tag, val reflect.Value, row int) {
	v := val.Interface()
	if tg.isEnum {
		key := fmt.Sprintf("%v", v)
		if mapped, ok := tg.Enum[key]; ok {
			v = mapped
		} else {
			v = "未知"
		}
	}
	g.grid[row-1][g.colIdx[tg.Column]-1] = v
}

// parseStyleJSON 把 v1 风格的样式 JSON(snake_case 键)转成 excelize/v2 的 *Style。
// v2 的 Style/Alignment/Font 结构体没有 json tag 且字段为驼峰, 直接 json.Unmarshal 会漏掉
// wrap_text / text_rotation / shrink_to_fit 这类下划线键, 故此处用带 json tag 的中间结构显式映射。
func parseStyleJSON(s string) (*excelize.Style, error) {
	var raw struct {
		Alignment *struct {
			Horizontal   string `json:"horizontal"`
			Vertical     string `json:"vertical"`
			WrapText     bool   `json:"wrap_text"`
			TextRotation int    `json:"text_rotation"`
			Indent       int    `json:"indent"`
			ShrinkToFit  bool   `json:"shrink_to_fit"`
		} `json:"alignment"`
		Font *struct {
			Bold   bool    `json:"bold"`
			Italic bool    `json:"italic"`
			Color  string  `json:"color"`
			Size   float64 `json:"size"`
			Family string  `json:"family"`
		} `json:"font"`
		Fill *struct {
			Type    string   `json:"type"`
			Pattern int      `json:"pattern"`
			Color   []string `json:"color"`
		} `json:"fill"`
	}
	if err := json.Unmarshal([]byte(s), &raw); err != nil {
		return nil, err
	}
	st := &excelize.Style{}
	if a := raw.Alignment; a != nil {
		st.Alignment = &excelize.Alignment{
			Horizontal:   a.Horizontal,
			Vertical:     a.Vertical,
			WrapText:     a.WrapText,
			TextRotation: a.TextRotation,
			Indent:       a.Indent,
			ShrinkToFit:  a.ShrinkToFit,
		}
	}
	if fo := raw.Font; fo != nil {
		st.Font = &excelize.Font{
			Bold:   fo.Bold,
			Italic: fo.Italic,
			Color:  fo.Color,
			Size:   fo.Size,
			Family: fo.Family,
		}
	}
	if fi := raw.Fill; fi != nil {
		st.Fill = excelize.Fill{
			Type:    fi.Type,
			Pattern: fi.Pattern,
			Color:   fi.Color,
		}
	}
	return st, nil
}
