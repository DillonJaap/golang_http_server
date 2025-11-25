package slices

import "slices"

func Flatten[T any](slice [][]T) []T {
	var newSlice []T
	for _, s := range slice {
		newSlice = slices.Concat(newSlice, s)
	}
	return newSlice
}
