// Description: Helper functions for validation and error handling.
package helper

import "reflect"

// IsNilOrEmpty checks if a value is nil or empty.
// It returns true if the value is nil, an empty string, or a zero value.
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
