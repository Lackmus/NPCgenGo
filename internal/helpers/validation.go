// Description: Helper functions for validation and error handling.
package helper

import "reflect"

func IsNilOrEmpty[T any](t T) bool {
	if reflect.ValueOf(t).IsZero() {
		return true
	}

	if str, ok := any(t).(string); ok && str == "" {
		return true
	}

	return false
}
