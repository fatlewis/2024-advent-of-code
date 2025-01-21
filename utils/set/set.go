package set

import (
	"slices"
)

func Union[T comparable](a []T, b []T) (result []T) {
	return slices.Concat(SetDifference(a, b), b)
}

func Intersection[T comparable](a []T, b []T) (result []T) {
	for _, aElem := range a {
		if slices.Contains(b, aElem) {
			result = append(result, aElem)
		}
	}
	return result
}

func SetDifference[T comparable](a []T, b []T) (result []T) {
	for _, aElem := range a {
		if !slices.Contains(b, aElem) {
			result = append(result, aElem)
		}
	}
	return result
}

func SymmetricDifference[T comparable](a []T, b []T) (result []T) {
	return slices.Concat(SetDifference(a, b), SetDifference(b, a))
}

func CartesianProduct[T comparable](a []T, b []T) (result [][2]T) {
	for _, aElem := range a {
		for _, bElem := range b {
			result = append(result, [2]T{aElem, bElem})
		}
	}
	return result
}

