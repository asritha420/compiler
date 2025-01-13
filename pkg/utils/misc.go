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

func Remove[T any](s []T, i int) []T {
	return append(s[:i], s[i+1:]...)
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

/*
Adds key and value from src to dst for any key dst does not have and it does not equal the ignore value. Returns the number of values added to dst.
*/
func AddToMapIgnore[K comparable, V any](src, dst map[K]V, ignore K) int {
	added := 0
	for k, v := range src {
		if k == ignore {
			continue
		}
		if _, ok := dst[k]; !ok {
			dst[k] = v
			added++
		}
	}
	return added
}

func HasSameKeys[K comparable, V any](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k := range m1 {
		if _, ok := m2[k]; !ok {
			return false
		}
	}

	return true
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

func AssertEqual[T comparable](t *testing.T, varName string, expected T, actual T) {
	if expected != actual {
		t.Fatalf("Expected %s=%v got %s=%v", varName, expected, varName, actual)
	}
}