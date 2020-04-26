package untils

func IsContain(v interface{}, s []string) bool {
	for _, value := range s {
		if v == value {
			return true
		}
	}

	return false
}
