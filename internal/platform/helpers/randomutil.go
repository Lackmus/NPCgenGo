// Description: This file contains helper functions for generating random values.
package helpers

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var (
	rndMu sync.Mutex
	rnd   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// SetSeed replaces the RNG with a new source seeded by `seed`.
// Useful for making tests deterministic.
func SetSeed(seed int64) {
	rndMu.Lock()
	defer rndMu.Unlock()
	rnd = rand.New(rand.NewSource(seed))
}

// SetRand injects a custom *rand.Rand. Use for advanced testing.
func SetRand(r *rand.Rand) {
	rndMu.Lock()
	defer rndMu.Unlock()
	if r != nil {
		rnd = r
	}
}

func GetRandomElement[T any](elements []T) T {
	if len(elements) == 0 {
		var zero T
		return zero
	}
	rndMu.Lock()
	idx := rnd.Intn(len(elements))
	rndMu.Unlock()
	return elements[idx]
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
	rndMu.Lock()
	n := rnd.Intn(max-min+1) + min
	rndMu.Unlock()
	return strconv.Itoa(n)
}

const (
	Random = "random"
)
