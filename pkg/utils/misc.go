package utils

import (
	"hash/fnv"
	"testing"
)

/*
Quickly remove an element at index i from array s.
*Note* There is no bounds checks and this will change the order of the array
*/
func FastRemove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func HashStr(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func HashArr[T Hashable](arr []T) int {
	sum := 0
	for _, elm := range arr {
		sum += elm.Hash()
	}
	return sum
}

func CompArrPtr[T Comparable[T]](arr1, arr2 []*T) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i, elm := range arr1 {
		if !(*elm).Equal(*arr2[i]) {
			return false
		}
	}
	return true
}

func AssertEqual[T comparable](t *testing.T, varName string, expected T, actual T) {
	if expected != actual {
		t.Fatalf("Expected %s=%v got %s=%v", varName, expected, varName, actual)
	}
}