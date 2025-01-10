package utils

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
