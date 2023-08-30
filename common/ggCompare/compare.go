package ggCompare

import "common/ggConv"

// Inter 同时存在于a和b的元素
func Inter[T comparable](a, b []T) []T {
	m := ggConv.Array2Map(b)
	var res []T
	for _, v := range a {
		if _, ok := m[v]; ok {
			res = append(res, v)
		}
	}
	return res
}

// Diff a中存在，b中不存在的元素
func Diff[T comparable](a, b []T) []T {
	m := ggConv.Array2Map(b)
	var res []T
	for _, v := range a {
		if _, ok := m[v]; !ok {
			res = append(res, v)
		}
	}
	return res
}

// Union a和b的并集
func Union[T comparable](a, b []T) []T {
	m := ggConv.Array2Map(a)
	for _, v := range b {
		if _, ok := m[v]; !ok {
			a = append(a, v)
		}
	}
	return a
}
