package helpers

func Filter(subarray []string, test func(string) bool) (ret []string) {
	for _, s := range subarray {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func Contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
