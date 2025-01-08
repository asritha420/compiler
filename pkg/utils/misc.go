package utils

import "testing"

/*
Quickly remove an element at index i from array s.
*Note* There is no bounds checks and this will change the order of the array
*/
func FastRemove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}


func AssertEqual[T comparable](t *testing.T, varName string, expected T, actual T) {
	if expected != actual {
		t.Fatalf("Expected %s=%v got %s=%v", varName, expected, varName, actual)
	}
}