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
