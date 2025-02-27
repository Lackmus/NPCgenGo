package helper

import "reflect"

// IsNilOrEmpty checks if a given value is nil (for pointers/interfaces) or empty (for strings, slices, maps).
func IsNilOrEmpty[T any](t T) bool {
	// Handle nil values (pointers, interfaces)
	if reflect.ValueOf(t).IsZero() {
		return true
	}

	// Handle empty string case
	if str, ok := any(t).(string); ok && str == "" {
		return true
	}

	return false
}
