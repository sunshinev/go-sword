package untils

func IsContain(v interface{}, s []string) bool {
	for _, value := range s {
		if v == value {
			return true
		}
	}

	return false
}

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
