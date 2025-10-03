package sets

import (
	"cmp"
)

// SortedIntersection returns the intersection of
// two sorted and dereplicated slices a and b.
func SortedIntersection[T cmp.Ordered](a, b []T) []T {
	var result []T
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch cmp.Compare(a[i], b[j]) {
		case 0:
			result = append(result, a[i])
			i++
			j++
		case 1:
			j++
		case -1:
			i++
		}
	}
	return result
}

// SortedUnion returns the union of
// two sorted and dereplicated slices a and b.
func SortedUnion[T cmp.Ordered](a, b []T) []T {
	var result []T
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch cmp.Compare(a[i], b[j]) {
		case 0:
			result = append(result, a[i])
			i++
			j++
		case 1:
			result = append(result, b[j])
			j++
		case -1:
			result = append(result, a[i])
			i++
		}
	}
	// Add remaining elements.
	result = append(result, a[i:]...)
	result = append(result, b[j:]...)
	return result
}

// SortedIntersectionLen returns the length of the intersection of
// two sorted and dereplicated slices a and b.
func SortedIntersectionLen[T cmp.Ordered](a, b []T) int {
	result := 0
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch cmp.Compare(a[i], b[j]) {
		case 0:
			i++
			j++
			result++
		case 1:
			j++
		case -1:
			i++
		}
	}
	return result
}

// SortedUnionLen returns the length of the union of
// two sorted and dereplicated slices a and b.
func SortedUnionLen[T cmp.Ordered](a, b []T) int {
	result := 0
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		result++
		switch cmp.Compare(a[i], b[j]) {
		case 0:
			i++
			j++
		case 1:
			j++
		case -1:
			i++
		}
	}
	// Add remaining elements.
	result += len(a) - i
	result += len(b) - j
	return result
}
