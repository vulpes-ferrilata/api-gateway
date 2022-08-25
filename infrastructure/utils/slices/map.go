package slices

import "github.com/pkg/errors"

func Map[T1 comparable, T2 comparable](mapper mapperFunc[T1, T2], slice []T1) ([]T2, error) {
	newSlice := make([]T2, 0)

	for _, element := range slice {
		result, err := mapper(element)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		newSlice = append(newSlice, result)
	}

	return newSlice, nil
}
