package slices

func Contains[T comparable](slice []T, elements ...T) bool {
	for sliceIdx := range slice {
		for elementsIdx := range elements {
			if slice[sliceIdx] == elements[elementsIdx] {
				return true
			}
		}
	}

	return false
}
