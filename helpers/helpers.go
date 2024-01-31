package helpers

func IsInArray[V comparable](val V, arr []V) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func RemoveIndex(s []interface{}, index int) []interface{} {
	return append(s[:index], s[index+1:]...)
}
