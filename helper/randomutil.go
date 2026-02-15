// Description: This file contains helper functions for generating random values.
package helper

import (
	"math/rand"
	"strconv"
)

func GetRandomElement[T any](elements []T) T {
	if len(elements) == 0 {
		var zero T
		return zero
	}
	return elements[rand.Intn(len(elements))]
}

// NewRandomMapKeySelector returns a function that, when called,
// returns a random key from the provided static map.
// The keys are computed and cached only once.
func NewRandomMapKeySelector[K comparable, V any](m map[K]V) func() K {
	// Cache the keys of the map.
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Return a closure that picks a random key from the cached keys.
	return func() K {
		return GetRandomElement(keys)
	}
}

func RandomInt(min, max int) string {
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

const (
	Random = "random"
)
