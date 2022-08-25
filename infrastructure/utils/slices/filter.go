package slices

func Filter[T comparable](predicate predicateFunc[T], slice []T) []T {
	newSlice := make([]T, 0)

	for _, element := range slice {
		if predicate(element) {
			newSlice = append(newSlice, element)
		}
	}

	return newSlice
}
