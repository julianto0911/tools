package lib_array

func IntPushDistinct(data []int, val int) []int {
	for i := range data {
		if data[i] == val {
			return data
		}
	}
	data = append(data, val)
	return data
}

func StrPushDistinct(data []string, val string) []string {
	for i := range data {
		if data[i] == val {
			return data
		}
	}
	data = append(data, val)
	return data
}
