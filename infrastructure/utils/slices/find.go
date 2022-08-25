package slices

func Find[T comparable](predicate predicateFunc[T], slice []T) (T, bool) {
	var empty T

	for _, element := range slice {
		if predicate(element) {
			return element, true
		}
	}

	return empty, false
}
