package slices

func Any[T comparable](predicate predicateFunc[T], slice []T) bool {
	for _, element := range slice {
		if predicate(element) {
			return true
		}
	}

	return false
}
