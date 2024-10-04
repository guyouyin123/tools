package client

import (
	"errors"
	"fmt"
	protocolEntry "github.com/withlin/canal-go/protocol/entry"
	"strings"
)

func logWhenOthers(entry protocolEntry.Entry) {
	if Debug {
		fmt.Println("--------logWhenOthers--------")
		fmt.Printf("%+v \n", entry)
		fmt.Println("-----------------------------")
	}
}

func FilterMap(header *protocolEntry.Header) bool {
	v, ok := SchemaNameMap[header.GetSchemaName()]
	if !v || !ok {
		return false
	}
	v, ok = TableMap[header.GetTableName()]
	if !v || !ok {
		return false
	}
	return true
}

func insertSql(header *protocolEntry.Header, rowData *protocolEntry.RowData) (string, error) {
	sqlStr, err := GetInsertSql(header, rowData)
	if err != nil {
		return "", err
	}
	return sqlStr, nil
}

func updateSql(header *protocolEntry.Header, rowData *protocolEntry.RowData) (string, error) {
	sqlStr, err := GetUpdateSql(header, rowData)
	if err != nil {
		return "", err
	}
	return sqlStr, nil
}

func deleteSql(header *protocolEntry.Header, rowData *protocolEntry.RowData) (string, error) {
	sqlStr, err := GetDeleteSql(header, rowData)
	if err != nil {
		return "", err
	}
	return sqlStr, nil
}

func GetInsertSql(header *protocolEntry.Header, rowData *protocolEntry.RowData) (string, error) {
	afterColumns := rowData.AfterColumns
	if len(afterColumns) == 0 {
		return "", errors.New("len(rowData.AfterColumns) == 0 ")
	}

	keys := ""
	values := ""
	for _, column := range afterColumns {
		if keys != "" {
			keys += ","
			values += ","
		}

		keys += "`" + column.GetName() + "`"
		values += GetValue(column)
	}
	tableName := "`" + header.GetSchemaName() + "`.`" + header.GetTableName() + "`"
	sqlInsert := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, keys, values)

	return sqlInsert, nil
}

func GetValue(column *protocolEntry.Column) string {
	if column.GetIsNull() {
		return "null"
	}

	// 怎么会出现单个 '\' ?
	return "'" + strings.ReplaceAll(strings.ReplaceAll(column.GetValue(), "\\", "\\\\"), "'", "\\'") + "'"
}
func GetUpdateSql(header *protocolEntry.Header, rowData *protocolEntry.RowData) (string, error) {
	beforeColumns := rowData.BeforeColumns
	afterColumns := rowData.AfterColumns
	if len(afterColumns) == 0 {
		return "", errors.New("len(rowData.AfterColumns) == 0 ")
	}

	sets := ""
	where := ""
	if HasKeyColumn(afterColumns) {
		for _, column := range afterColumns {
			if column.GetIsKey() {
				if where != "" {
					where += " and "
				}

				condition := "="
				if column.GetIsNull() {
					condition = " is "
				}
				where += "`" + column.GetName() + "`" + condition + GetValue(column)
				continue // 主键 where clause
			}

			if column.GetUpdated() {
				if sets != "" {
					sets += ","
				}
				sets += "`" + column.GetName() + "`" + "=" + GetValue(column)
			}
		}
	} else {
		for _, column := range afterColumns {
			if sets != "" {
				sets += ","
			}
			sets += "`" + column.GetName() + "`" + "=" + GetValue(column)
		}

		for _, column := range beforeColumns {
			if where != "" {
				where += " and "
			}

			condition := "="
			if column.GetIsNull() {
				condition = " is "
			}
			where += "`" + column.GetName() + "`" + condition + GetValue(column)
		}
	}
	tableName := "`" + header.GetSchemaName() + "`.`" + header.GetTableName() + "`"
	sqlUpdate := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, sets, where)
	return sqlUpdate, nil
}

// 检查是否有主键列
func HasKeyColumn(cols []*protocolEntry.Column) bool {
	for _, col := range cols {
		if col.GetIsKey() {
			return true
		}
	}

	return false
}

func GetDeleteSql(header *protocolEntry.Header, rowData *protocolEntry.RowData) (string, error) {
	beforeColumns := rowData.BeforeColumns
	if len(beforeColumns) == 0 {
		return "", errors.New("len(rowData.AfterColumns) == 0 ")
	}

	where := ""
	if HasKeyColumn(beforeColumns) {
		for _, column := range beforeColumns {
			if column.GetIsKey() {
				if where != "" {
					where += " and "
				}

				condition := "="
				if column.GetIsNull() {
					condition = " is "
				}
				where += "`" + column.GetName() + "`" + condition + GetValue(column)
			}
		}
	} else {
		for _, column := range beforeColumns {
			if where != "" {
				where += " and "
			}

			condition := "="
			if column.GetIsNull() {
				condition = " is "
			}
			where += "`" + column.GetName() + "`" + condition + GetValue(column)
		}
	}
	tableName := "`" + header.GetSchemaName() + "`.`" + header.GetTableName() + "`"
	sqlDelete := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, where)

	return sqlDelete, nil
}
