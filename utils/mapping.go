package utils

func MapSlice[T any, R any](s []T, f func(v T, i int) R) []R {
	result := make([]R, len(s))
	for i, v := range s {
		result[i] = f(v, i)
	}
	return result
}

func PruneDuplicates[T comparable](s []T) []T {
	result := make([]T, 0, len(s))
	seen := make(map[T]bool)
	for _, v := range s {
		if !seen[v] {
			result = append(result, v)
			seen[v] = true
		}
	}
	return result
}

func FilterSlice[T any](s []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func FilterEmpty[T string](s []T) []T {
	return FilterSlice(s, func(v T) bool {
		return v != ""
	})
}
