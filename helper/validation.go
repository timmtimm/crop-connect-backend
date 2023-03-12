package helper

func CheckStringOnArray(array []string, check string) bool {
	for _, v := range array {
		if v == check {
			return true
		}
	}
	return false
}
