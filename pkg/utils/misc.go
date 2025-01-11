package utils

import (
	"fmt"
	"hash/fnv"
	"strings"
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

/*
Adds key and value from src to dst for any key dst does not have. Returns the number of values added to dst.
*/
func AddToMap[K comparable, V any](src, dst map[K]V) int {
	added := 0
	for s, v := range src {
		if _, ok := dst[s]; !ok {
			dst[s] = v
			added++
		}
	}
	return added
}

func MapToSetString[K comparable, V any](set map[K]V) string {
	strs := make([]string, len(set))
	i := 0
	for val := range set {
		strs[i] = fmt.Sprint(val)
		i++
	}
	return fmt.Sprintf("{%s}", strings.Join(strs, ","))
}

func MapHasSameKeys[K comparable, V any](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}

	for key := range m1 {
		if _, ok := m2[key]; !ok {
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