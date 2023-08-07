package conv

func Array2Map[T comparable](array []T) map[T]struct{} {
	m := make(map[T]struct{}, len(array))
	for _, v := range array {
		m[v] = struct{}{}
	}
	return m
}
