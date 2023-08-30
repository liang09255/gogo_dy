package ggConv

func Array2Map[T comparable](array []T) map[T]struct{} {
	m := make(map[T]struct{}, len(array))
	for _, v := range array {
		m[v] = struct{}{}
	}
	return m
}

func Array2BoolMap[T comparable](array []T) map[T]bool {
	m := make(map[T]bool, len(array))
	for _, v := range array {
		m[v] = true
	}
	return m
}

// Map2Array 只保留key的信息
func Map2Array[T comparable](m map[T]any) []T {
	array := make([]T, 0, len(m))
	for k, _ := range m {
		array = append(array, k)
	}
	return array
}
