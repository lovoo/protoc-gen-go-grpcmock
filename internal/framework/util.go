package framework

import (
	"strings"
)

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }

func substringAfter(value, sep string) string {
	pos := strings.Index(value, sep)
	if pos == -1 {
		return ""
	}
	start := pos + len(sep)
	if start >= len(value) {
		return ""
	}
	return value[start:]
}

func mapSlice[T any, S any](a []T, f func(T) S) []S {
	n := make([]S, len(a))
	for i, e := range a {
		n[i] = f(e)
	}
	return n
}
