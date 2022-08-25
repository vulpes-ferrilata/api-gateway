package slices

type mapperFunc[T1 comparable, T2 comparable] func(object T1) (T2, error)
