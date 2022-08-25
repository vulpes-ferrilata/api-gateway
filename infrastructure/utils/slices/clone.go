package slices

func Clone[T comparable](slice []T) []T {
	newSlice := make([]T, len(slice))
	copy(newSlice, slice)
	return newSlice
}
