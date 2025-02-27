// Description: Helper functions for copying maps and slices
package helper

// CopySlice : Copy a slice
// Note: This function only works for slices of basic types or slices of maps or slices.
// It does not work for slices of structs or slices of slices.
func CopyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// DeepCopySlice : Deep copy a slice
// Note: This function only works for slices of basic types or slices of maps or slices.
// It does not work for slices of structs or slices of slices.
func DeepCopyMap[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		switch val := any(v).(type) {
		case map[K]V:
			dst[k] = any(DeepCopyMap(val)).(V) // Recursive map copy
		case []V:
			newSlice := make([]V, len(val))
			copy(newSlice, val) // Copy slice contents
			dst[k] = any(newSlice).(V)
		default:
			dst[k] = v
		}
	}
	return dst
}
