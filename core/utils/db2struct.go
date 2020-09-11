package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sunshinev/go-sword/config"
)

type Db2struct struct {
}

type RowItem struct {
	ColumnName    string
	ColumnKey     string // PRI
	DataType      string
	IsNullable    string // NO
	ColumnComment string
}

func (ds Db2struct) Convert(table string) (*[]RowItem, error) {
	rows, err := config.GlobalConfig.DbConn.Query("SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE,COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND table_name = ?", config.GlobalConfig.DatabaseSet.Database, table)
	if err != nil {
		return nil, err
	}

	if rows != nil {
		defer func() {
			_ = rows.Close()
		}()
	} else {
		return nil, errors.New("No results returned for table")
	}

	columns := []RowItem{}
	for rows.Next() {
		var column, columnKey, dataType, nullable, comment string
		err := rows.Scan(&column, &columnKey, &dataType, &nullable, &comment)
		if err != nil {
			return nil, err
		}
		columns = append(columns, RowItem{
			ColumnName:    column,
			ColumnKey:     columnKey,
			DataType:      dataType,
			IsNullable:    nullable,
			ColumnComment: comment,
		})
	}

	return &columns, nil
}

// 得到整体的结构体文件
func (s Db2struct) FetchWholeStructFile(packageName, structName, tableName string, columns *[]RowItem) string {
	structBody := s.GenerateStruct(structName, columns)
	importStr := ""
	if strings.Contains(structBody, "time.Time") {
		importStr = "import \"time\""
	}
	return fmt.Sprintf("%s\n\n%s\n%s\n\n%s", s.genPackage(packageName), importStr, structBody, s.genTableNameFunc(structName, tableName))
}

func (s Db2struct) GenerateStruct(structName string, columns *[]RowItem) string {
	itemList := []string{}
	for _, row := range *columns {
		rowStr := fmt.Sprintf("\t%s %s `%s %s`", strings.Replace(strings.Title(strings.Replace(row.ColumnName, "_", " ", -1)), " ", "", -1), s.mappingType(row.DataType), s.mappingGormTag(&row), s.mappingJsonTag(&row))
		itemList = append(itemList, rowStr)
	}

	st := fmt.Sprintf("type %s struct {\n%s\n}", structName, strings.Join(itemList, "\n"))
	return st
}

// Mysql 类型映射到 Go
func (s Db2struct) mappingType(fieldType string) string {
	switch fieldType {
	case "char", "varchar", "tinytext", "text", "blob", "mediumtext", "mediumblob", "longblob", "longtext", "enum", "json":
		return "string"
	case "date", "datetime", "timestamp", "time", "year":
		return "time.Time"
	case "tinyint", "smallint", "mediumint", "int", "bigint":
		return "int32"
	case "float", "double", "decimal":
		return "float64"
	default:
		return "string"
	}
}

func (s Db2struct) mappingJsonTag(row *RowItem) string {
	return fmt.Sprintf("json:\"%s\"", row.ColumnName)
}

func (s Db2struct) mappingGormTag(row *RowItem) string {
	isNullable := ""
	if row.IsNullable == "YES" {
		isNullable = "not null"
	}
	isPk := ""
	if row.ColumnKey == "PRI" {
		isPk = "primaryKey"
	}

	tags := []string{"Column:" + row.ColumnName, isNullable, isPk}
	newTags := []string{}
	for _, v := range tags {
		if v != "" {
			newTags = append(newTags, v)
		}
	}

	return fmt.Sprintf("gorm:\"%s\"", strings.Join(newTags, ";"))
}

func (s Db2struct) genTableNameFunc(structName, tableName string) string {
	return `func(` + strings.ToLower(structName) + ` *` + structName + `)TableName() string {
	return "` + tableName + `"
}`
}

func (s Db2struct) genPackage(packageName string) string {
	return fmt.Sprintf("package %s", packageName)
}
