package utils

// This function checks if a value exists in an array or not
func CheckFileFormat(a [3]string, val string) bool {
	for _, v := range a {
		if v == val {
			return true
		}
	}
	return false
}
