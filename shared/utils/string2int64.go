package utils

import "strconv"

func String2Int(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func Int2String(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StringArray2Int64Array(s []string) []int64 {
	var result []int64
	for _, v := range s {
		result = append(result, String2Int(v))
	}
	return result
}

func Int64Array2StringArray(i []int64) []string {
	var result []string
	for _, v := range i {
		result = append(result, Int2String(v))
	}
	return result
}
