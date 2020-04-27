package untils

import "sort"

func IsContain(v interface{}, s []string) bool {
	for _, value := range s {
		if v == value {
			return true
		}
	}

	return false
}

// Convert mysql field type to number or string
// Because when web post json to background ,the type may not match struct
func ConvertFieldsType2Js(mysqlType string) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		return "number"
	case "bigint":
		return "number"
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		return "string"
	case "date", "datetime", "time", "timestamp":
		return "string"
	case "decimal", "double":
		return "number"
	case "float":
		return "number"
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return "string"
	}

	return "string"
}

// Because db2struct output fields is sort by `sort.strings`
// so we will put `id` at the top,and `create_at`,`updated_at` at bottom
func ResortMySQLFields(fields *[]string) []string {
	newFields := []string{"id"}

	sort.Strings(*fields)

	for _, f := range *fields {
		if f == "id" || f == "created_at" || f == "updated_at" {
			continue
		}
		newFields = append(newFields, f)
	}

	newFields = append(newFields, "created_at")
	newFields = append(newFields, "updated_at")

	return newFields
}
