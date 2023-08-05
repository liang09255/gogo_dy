package service

import "strconv"

// GetInt64Bystring string类型转int64
func GetInt64Bystring(str string) int64 {
	res, _ := strconv.Atoi(str)
	return int64(res)
}
