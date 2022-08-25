package slices

type predicateFunc[T comparable] func(object T) bool
