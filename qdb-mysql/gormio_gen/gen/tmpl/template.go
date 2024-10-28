package tmpl

import (
	"bytes"
	"fmt"
	"github.com/guyouyin123/tools/qdb-mysql/gormio_gen/gen/conf"
	"os"
	"strings"
	"text/template"
	"unicode"
)

type tempConf struct {
	TableName  string
	StructName string
	FieldName  string
	FileType   string
}

func RunTemplate(basePath string) {
	if !conf.IsTemplate {
		return
	}
	for tableName, fieldNameList := range conf.TemplateConf {
		path := fmt.Sprintf("%s/%s.gen.go", basePath, tableName)
		structNameZ := ToUpperGorm(tableName)
		tableNameZ := ToLowerGorm(tableName)
		contentList := make([]string, 0)
		for fieldName, fileType := range fieldNameList {
			fieldNameZ := ToUpperGorm(fieldName)
			infosByIdsTemp := GetTmpl(tableNameZ, fieldNameZ, structNameZ, fileType)
			contentList = append(contentList, infosByIdsTemp)
		}
		templateWrite(path, contentList)
	}
}

func GetTmpl(tableName, fieldName, structName, fileType string) string {
	GetInfosByIdsTmplStr := `
// GetInfosBy{{.FieldName}}s DO NOT EDIT.
func (temp *{{.TableName}}) GetInfosBy{{.FieldName}}s(ids []{{.FileType}}) ([]*model.{{.StructName}}, error) {
	infos,err:= temp.Where(temp.{{.FieldName}}.In(ids...)).Find()
	if err!=nil{
		return nil, err
	}
	return infos, nil
}

// GetInfoBy{{.FieldName}} DO NOT EDIT.
func (temp *{{.TableName}}) GetInfoBy{{.FieldName}}(id {{.FileType}}) (*model.{{.StructName}}, error) {
	info,err:= temp.Where(temp.{{.FieldName}}.Eq(id)).First()
	if err!=nil{
		return nil, err
	}
	return info, nil
}

// GetInfosMapBy{{.FieldName}}s DO NOT EDIT.
func (temp *{{.TableName}}) GetInfosMapBy{{.FieldName}}s(ids []{{.FileType}}) (map[{{.FileType}}]*model.{{.StructName}}, error) {
	infos,err:= temp.Where(temp.{{.FieldName}}.In(ids...)).Find()
	if err!=nil{
		return nil, err
	}
	m:=map[{{.FileType}}]*model.{{.StructName}}{}
	for _,v:=range infos{
		m[v.{{.FieldName}}]=v
	}
	return m, nil
}
`
	data := tempConf{
		TableName:  tableName,
		StructName: structName,
		FieldName:  fieldName,
		FileType:   fileType,
	}
	buf := &bytes.Buffer{}
	tmpl, err := template.New(tableName).Parse(GetInfosByIdsTmplStr)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func templateWrite(path string, contentList []string) {
	// 打开文件，如果文件不存在则创建它
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// 追加内容
	for _, content := range contentList {
		if _, err := file.WriteString(content); err != nil {
			panic(err)
		}
	}
}

// ToUpperGorm 大驼峰Gorm.io
func ToUpperGorm(s string) string {
	//gorm.io/gen 生成模型时,特定的字段名转换为全大写
	upperMap := map[string]string{
		"id":   "ID",
		"uuid": "UUID",
		"guid": "GUID",
	}
	words := strings.Split(s, "_")
	for i, w := range words {
		f, ok := upperMap[w]
		if ok {
			words[i] = f
		} else if w != "" {
			words[i] = strings.ToUpper(string(unicode.ToUpper(rune(w[0])))) + w[1:]
		}
	}
	res := strings.Join(words, "")
	//如果以s结尾，过滤--gorm.io特性。eg:表名users,生成的结构体为user
	if strings.HasSuffix(res, "s") {
		return res[:len(res)-1]
	}
	return res
}

// ToLowerGorm 小驼峰Gorm.io
func ToLowerGorm(s string) string {
	words := strings.Split(s, "_")
	for i, w := range words {
		if i == 0 {
			continue
		}
		if w != "" {
			words[i] = strings.ToUpper(string(unicode.ToUpper(rune(w[0])))) + w[1:]
		}
	}
	res := strings.Join(words, "")
	//如果以s结尾，过滤--gorm.io特性。eg:表名users,生成的结构体为user
	if strings.HasSuffix(res, "s") {
		return res[:len(res)-1]
	}
	return res
}
