package utils

type Slice[T any] struct {
	Values []T
}

func (slice *Slice[T]) Append(newValues ...T) *Slice[T] {
	slice.Values = append(slice.Values, newValues...)
	return slice
}

func (slice *Slice[T]) AppendIf(condition bool, newValues ...T) *Slice[T] {
	if condition {
		slice.Values = append(slice.Values, newValues...)
	}
	return slice
}

func Prepend[T any](slice []T, elems ...T) []T {
	return append(append([]T{}, elems...), slice...)
}

// remove if not used
func Map[T any, Q any](slice []T, f func(T) Q) []Q {
	res := make([]Q, len(slice))
	for i := 0; i < len(slice); i++ {
		res[i] = f(slice[i])
	}
	return res
}

func Filter[T any](slice []T, fn func(T) bool) []T {
	var res []T
	for _, x := range slice {
		if fn(x) {
			res = append(res, x)
		}
	}
	return res
}

func FirstNonZero[T comparable](values ...T) T {
	var zero T
	for _, val := range values {
		if val != zero {
			return val
		}
	}
	return zero
}

func Group[T any, Q comparable](slice []T, fn func(T) Q) [][]T {
	var all [][]T
	var group []T
	for _, x := range slice {
		key := fn(x)
		if len(group) > 0 {
			if fn(group[0]) == key {
				group = append(group, x)
			} else {
				all = append(all, group)
				group = nil
			}
		} else {
			group = append(group, x)
		}
	}
	if len(group) > 0 {
		all = append(all, group)
	}
	return all
}
